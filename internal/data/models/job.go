package models

import "time"

type JobStatus string

const JobStatusCreated JobStatus = "created"
const JobStatusPending JobStatus = "graded"
const JobStatusIgnored JobStatus = "ignored"
const JobStatusApplied JobStatus = "applied"

type Job struct {
	Id             string
	Title          string
	Url            string
	Description    string
	DatePosted     time.Time
	Company        string
	Location       string
	NumApplicants  string
	Status         JobStatus
	Grade          *int
	GradeReasoning *string
	Note           string
	IsInvalid      bool
}
