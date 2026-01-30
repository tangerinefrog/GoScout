package main

import (
	"context"
	"fmt"
	"job-scraper/internal/data"
	"job-scraper/internal/data/repositories"
	"job-scraper/internal/handlers"
	"job-scraper/internal/services/scheduler"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("could not load .env file: %w", err)
	}

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		return fmt.Errorf("server address is not defined in the .env file")
	}

	db, err := data.Init()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = initConfig(ctx, db)
	if err != nil {
		return fmt.Errorf("config init error: %w", err)
	}

	srv := configureServer(addr, db)

	go func() {
		log.Printf("server is listening on '%s'...\n", addr)
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	startBackgroundJobs(ctx, db)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	srv.Shutdown(shutdownCtx)

	return nil
}

func configureServer(addr string, db *data.DB) *http.Server {
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	h := handlers.NewHandler(db)
	h.SetupRoutes(router)

	return &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

func initConfig(ctx context.Context, db *data.DB) error {
	configRepo := repositories.NewConfigRepo(db)
	err := configRepo.Init(ctx)

	if err != nil {
		return err
	}

	return nil
}

func startBackgroundJobs(ctx context.Context, db *data.DB) {
	periodHour := 1 * time.Hour
	go scheduler.ScrapeRecurring(ctx, periodHour, db)
}
