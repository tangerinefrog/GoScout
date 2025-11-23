package main

import (
	"job-scraper/internal/data"
	"job-scraper/internal/services/scraper"
	"log"
	"time"

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

	s := scraper.NewScraper(db)
	_, err = s.ScrapeLinkedInJobs("Golang", 1*time.Hour)
	if err != nil {
		log.Fatal(err)
	}
}
