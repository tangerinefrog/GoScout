package scheduler

import (
	"context"
	"job-scraper/internal/data"
	"job-scraper/internal/services/scraper"
	"log"
	"time"
)

func ScrapeRecurring(ctx context.Context, d time.Duration, db *data.DB, keyword string, filterKeywords []string) {
	ticker := time.NewTicker(d)
	defer ticker.Stop()
	s := scraper.NewScraper(db)

	for {
		tCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)

		start := time.Now()
		log.Println("Scraping process started")
		_, err := s.ScrapeLinkedInJobs(tCtx, "go golang", nil, d)
		if err != nil {
			log.Printf("recurrent scraping encountered an error: %v", err)
			continue
		}

		log.Printf("Scraping process ended, run time: %v", time.Since(start))

		cancel()

		<-ticker.C
	}
}
