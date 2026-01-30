package main

import (
	"context"
	"job-scraper/internal/data"
	"job-scraper/internal/data/repositories"
	"job-scraper/internal/handlers"
	"job-scraper/internal/services/scheduler"
	"log"
	"os"
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

	log.Printf("Server is listening on ':%s'...\n", addr)
	err = r.Run(addr)
	log.Printf("Server error: %v\n", err)
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
