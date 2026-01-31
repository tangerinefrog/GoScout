package scheduler

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/tangerinefrog/GoScout/internal/data"
	"github.com/tangerinefrog/GoScout/internal/data/repositories"
	"github.com/tangerinefrog/GoScout/internal/services/scraper"
)

type scrapingConfig struct {
	searchQuery    string
	filterKeywords []string
}

func ScrapeRecurring(ctx context.Context, period time.Duration, db *data.DB) {
	s := scraper.NewScraper(db)

	err := runScrape(ctx, db, s, period)
	if err != nil {
		log.Printf("Error occured during recurrent scraping: %v", err)
	}

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err = runScrape(ctx, db, s, period)
			if err != nil {
				log.Printf("Error occured during recurrent scraping: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func runScrape(ctx context.Context, db *data.DB, s *scraper.Scraper, timeWindow time.Duration) error {
	cfg, err := getConfig(ctx, db)
	if err != nil {
		return fmt.Errorf("failed to get config: %w", err)
	}

	scrapeCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)
	defer cancel()

	start := time.Now()
	log.Println("Scraping process started")
	err = s.ScrapeLinkedInJobs(scrapeCtx, cfg.searchQuery, cfg.filterKeywords, timeWindow)
	if err != nil {
		return err
	}

	log.Printf("Scraping process ended, run time: %v", time.Since(start))
	return nil
}

func getConfig(ctx context.Context, db *data.DB) (scrapingConfig, error) {
	configRepo := repositories.NewConfigRepo(db)

	config, err := configRepo.Get(ctx)
	if err != nil {
		return scrapingConfig{}, err
	}

	periodHours := config.SearchPeriodHours
	if periodHours <= 0 {
		return scrapingConfig{}, fmt.Errorf("search period value is set to an incorrect value in the config")
	}

	searchQuery := strings.TrimSpace(config.SearchQuery)
	if searchQuery == "" {
		return scrapingConfig{}, fmt.Errorf("search query value is not set in the config")
	}
	searchFilter := strings.TrimSpace(config.SearchFilter)
	filterKeywords := strings.Split(searchFilter, ",")

	return scrapingConfig{
		searchQuery:    searchQuery,
		filterKeywords: filterKeywords,
	}, nil
}
