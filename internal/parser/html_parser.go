package parser

import (
	"bytes"
	"job-scraper/internal/models"
	"strings"

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

func ParseJob(body []byte, id string) (models.JobPosition, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return models.JobPosition{}, err
	}

	job := models.JobPosition{}

	job.ID = id
	job.Title = findAndTrimText(doc, ".top-card-layout__title")
	job.TimeAgo = findAndTrimText(doc, ".posted-time-ago__text")
	job.CompanyName = findAndTrimText(doc, ".topcard__org-name-link")
	job.LocationName = findAndTrimText(doc, ".topcard__flavor.topcard__flavor--bullet")
	job.Description = findAndTrimHtml(doc, ".show-more-less-html__markup")
	applicantsText := findAndTrimText(doc, ".num-applicants__caption")
	job.PageUrl, _ = doc.Find(".topcard__link").Attr("href")

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
