package main

import (
	"context"
	"job-scraper/internal/data"
	"job-scraper/internal/handlers"
	"job-scraper/internal/services/scheduler"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("could not load .env file: %v", err)
	}

	db, err := data.Init()
	if err != nil {
		log.Fatal(err)
	}

	r := gin.Default()
	r.Use(cors.Default())

	h := handlers.NewHandler(db)
	h.SetupRoutes(r)

	addr := os.Getenv("SRV_ADDR")
	if addr == "" {
		log.Fatalf("Server address is not defined in the .env file")
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	startRecurrentJobs(ctx, db)

	log.Printf("Server is listening on ':%s'...\n", addr)
	err = r.Run(addr)
	log.Printf("Server error: %v\n", err)
}

func startRecurrentJobs(ctx context.Context, db *data.DB) {
	v := os.Getenv("SCRAPING_INTERVAL_HOURS")
	intervalHours, err := strconv.Atoi(v)
	if err != nil {
		log.Printf("Scraping interval is not defined in the .env file, default interval of 1 hour is set")
	}

	d := 1 * time.Hour
	if intervalHours > 0 {
		d = time.Duration(intervalHours) * time.Hour
	}

	keywords := os.Getenv("SCRAPING_KEYWORDS")
	filterBy := os.Getenv("SCRAPING_FILTER_BY")
	if keywords == "" {
		log.Print("Scraping keywords are not defined in the .env file, recurrent job is not started")
	} else {
		filterKeywords := strings.Split(filterBy, ",")
		go scheduler.ScrapeRecurring(ctx, d, db, keywords, filterKeywords)
	}
}
