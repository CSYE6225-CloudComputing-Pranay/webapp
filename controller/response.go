package controller

import "time"

type AccountResponse struct {
	ID             string    `json:"id"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	Password       string    `json:"password"`
	Email          string    `json:"email"`
	AccountCreated time.Time `json:"account_created"`
	AccountUpdated time.Time `json:"account_updated"`
}

type AssignmentResponse struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Points            int       `json:"points"`
	NumOfAttempts     int       `json:"num_of_attempts"`
	Deadline          time.Time `json:"deadline"`
	AssignmentCreated time.Time `json:"assignment_created"`
	AssignmentUpdated time.Time `json:"assignment_updated"`
}

type SubmissionResponse struct {
	ID                string    `json:"id"`
	AssignmentID      string    `json:"assignment_id"`
	SubmissionURL     string    `json:"submission_url"`
	SubmissionDate    time.Time `json:"submission_date"`
	SubmissionUpdated time.Time `json:"submission_updated"`
}
