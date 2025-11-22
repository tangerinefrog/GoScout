package main

import (
	"fmt"
	"job-scraper/internal/models"
	"job-scraper/internal/scraper"
	"log"
)

func main() {
	positions, err := scraper.ScrapeLinkedInJobs("Golang Developer")
	if err != nil {
		log.Fatal(err)
	}

	for _, v := range positions {
		printPosition(v)
	}

}

func printPosition(job models.JobPosition) {
	fmt.Printf("%s at '%s' \n\nUrl: %s\nLocation: %s\nPosted: %s\nApplicants: %s\n\n",
		job.Title, job.CompanyName, job.PageUrl, job.LocationName, job.TimeAgo, job.NumApplicants)
}
