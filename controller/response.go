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
