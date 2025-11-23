package models

import "time"

type Job struct {
	Id            string
	Title         string
	Url           string
	Description   string
	DatePosted    time.Time
	Company       string
	Location      string
	NumApplicants string
	Status        string
}
