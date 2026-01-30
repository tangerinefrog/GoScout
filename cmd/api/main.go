package main

import (
	"context"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Could not load .env file: %v", err)
	}

	db, err := data.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = initConfig(ctx, db)
	if err != nil {
		log.Fatalf("Config init error: %v", err)
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(cors.Default())

	h := handlers.NewHandler(db)
	h.SetupRoutes(r)

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		log.Fatalf("Server address is not defined in the .env file")
	}

	go startRecurrentJobs(ctx, db)

	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		log.Printf("Server is listening on '%s'...\n", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	srv.Shutdown(shutdownCtx)
}

func initConfig(ctx context.Context, db *data.DB) error {
	configRepo := repositories.NewConfigRepo(db)
	err := configRepo.Init(ctx)

	if err != nil {
		return err
	}

	return nil
}

func startRecurrentJobs(ctx context.Context, db *data.DB) {
	periodHour := 1 * time.Hour
	err := scheduler.ScrapeRecurring(ctx, periodHour, db)
	if err != nil {
		log.Printf("Recurrent scraping error: %v", err)
	}
}
