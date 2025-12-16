package scheduler

import (
	"context"
	"errors"
	"job-scraper/internal/data"
	"job-scraper/internal/data/repositories"
	"job-scraper/internal/services/scraper"
	"log"
	"strings"
	"time"
)

type scrapingConfig struct {
	searchQuery    string
	filterKeywords []string
}

func ScrapeRecurring(ctx context.Context, period time.Duration, db *data.DB) error {
	cfg, err := getConfig(ctx, db)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	s := scraper.NewScraper(db)

	for {
		tCtx, cancel := context.WithTimeout(ctx, 10*time.Minute)

		start := time.Now()
		log.Println("Scraping process started")
		_, err := s.ScrapeLinkedInJobs(tCtx, cfg.searchQuery, cfg.filterKeywords, period)
		if err != nil {
			log.Printf("recurrent scraping encountered an error: %v", err)
			continue
		}

		log.Printf("Scraping process ended, run time: %v", time.Since(start))

		cancel()

		<-ticker.C

		tmp, err := getConfig(ctx, db)
		if err != nil {
			log.Printf("Error occured while getting configuration for recurrent scraping: %v", err)
		} else {
			cfg = tmp
		}
	}
}

func getConfig(ctx context.Context, db *data.DB) (scrapingConfig, error) {
	configRepo := repositories.NewConfigRepo(db)

	config, err := configRepo.Get(ctx)
	if err != nil {
		return scrapingConfig{}, err
	}

	periodHours := config.SearchPeriodHours
	if periodHours <= 0 {
		return scrapingConfig{}, errors.New("search period value is set to an incorrect value in the config")
	}

	searchQuery := strings.TrimSpace(config.SearchQuery)
	if searchQuery == "" {
		return scrapingConfig{}, errors.New("search query value is not set in the config")
	}
	searchFilter := strings.TrimSpace(config.SearchFilter)
	filterKeywords := strings.Split(searchFilter, ",")

	return scrapingConfig{
		searchQuery:    searchQuery,
		filterKeywords: filterKeywords,
	}, nil
}
