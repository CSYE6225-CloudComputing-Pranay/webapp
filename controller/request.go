package controller

import "time"

type AccountRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type AssignmentRequest struct {
	ID                string    `json:"id"`
	Name              string    `json:"name" binding:"required"`
	Points            int       `json:"points" binding:"required,max=10,min=1"`
	NumOfAttempts     int       `json:"num_of_attempts" binding:"required,max=100,min=1"`
	Deadline          time.Time `json:"deadline" binding:"required"`
	AssignmentCreated time.Time `json:"assignment_created"`
	AssignmentUpdated time.Time `json:"assignment_updated"`
}
