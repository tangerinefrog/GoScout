package parser

import (
	"bytes"
	"errors"
	"fmt"
	"job-scraper/internal/data/models"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func ParseIdsFromSearch(body []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res := make([]string, 0, 10)

	doc.Find("body>li .base-card").Each(func(i int, s *goquery.Selection) {
		attrVal, ok := s.Attr("data-entity-urn")
		if ok {
			parts := strings.Split(attrVal, ":")
			if len(parts) >= 4 {
				res = append(res, parts[3])
			}

		}
	})

	return res, nil
}

func ParseJob(body []byte, id string) (models.Job, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return models.Job{}, err
	}

	job := models.Job{}

	job.Id = id
	job.Title = findAndTrimText(doc, ".top-card-layout__title")
	job.Company = findAndTrimText(doc, ".topcard__org-name-link")
	job.Location = findAndTrimText(doc, ".topcard__flavor.topcard__flavor--bullet")
	job.Description = findAndTrimHtml(doc, ".show-more-less-html__markup")
	job.Url, _ = doc.Find(".topcard__link").Attr("href")
	
	timeAgo := findAndTrimText(doc, ".posted-time-ago__text")
	date, err := extractDate(timeAgo)
	if err == nil {
		job.DatePosted = date
	}

	applicantsText := findAndTrimText(doc, ".num-applicants__caption")
	job.NumApplicants = getApplicants(applicantsText)

	return job, nil
}

func findAndTrimHtml(doc *goquery.Document, selector string) string {
	html, _ := doc.Find(selector).Html()
	return strings.TrimSpace(html)
}

func findAndTrimText(doc *goquery.Document, selector string) string {
	return strings.TrimSpace(doc.Find(selector).Text())
}

func getApplicants(s string) string {
	if s != "" {
		parts := strings.Split(s, " ")
		if parts[0] == "Be" {
			return "<" + parts[4]
		} else {
			return parts[0]
		}
	}

	return ""
}

func extractDate(timeAgo string) (time.Time, error) {
	timeAgo = strings.TrimSpace(timeAgo)
	if timeAgo == "" {
		return time.Time{}, errors.New("could not extract date from an empty string")
	}

	parts := strings.Split(timeAgo, " ")
	if len(parts) < 2 {
		return time.Time{}, fmt.Errorf("could not parse date from '%s'", timeAgo)
	}

	num, err := strconv.Atoi(parts[0])
	if len(parts) < 2 || err != nil {
		return time.Time{}, fmt.Errorf("could not parse date from '%s'", timeAgo)
	}

	name := strings.TrimRight(parts[1], "s")
	durations := map[string]time.Duration{
		"minute": time.Minute,
		"hour":   time.Hour,
		"day":    time.Hour * 24,
		"week":   time.Hour * 24 * 7,
		"month":  time.Hour * 24 * 30,
	}

	if dur, ok := durations[name]; ok {
		date := time.Now().Add(-1 * time.Duration(num) * dur)
		return date, nil
	}

	return time.Time{}, fmt.Errorf("could not parse date from '%s'; unknown duration '%s'", timeAgo, name)

}
