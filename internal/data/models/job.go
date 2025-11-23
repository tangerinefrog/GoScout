package models

import "time"

type jobStatus string

const JobStatusCreated jobStatus = "created"
const JobStatusPending jobStatus = "pending"
const JobStatusIgnored jobStatus = "ignored"
const JobStatusApplied jobStatus = "applied"

type Job struct {
	Id            string
	Title         string
	Url           string
	Description   string
	DatePosted    time.Time
	Company       string
	Location      string
	NumApplicants string
	Status        jobStatus
}
