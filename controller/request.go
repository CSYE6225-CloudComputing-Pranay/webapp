package controller

import "time"

type AccountRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

type AssignmentRequest struct {
	Name              string    `json:"name" binding:"required"`
	Points            int       `json:"points" validate:"max=100,min=1" binding:"required"`
	NumOfAttempts     int       `json:"num_of_attempts" validate:"max=10,min=1" binding:"required"`
	Deadline          time.Time `json:"deadline" binding:"required"`
	AssignmentCreated time.Time `json:"assignment_created"`
	AssignmentUpdated time.Time `json:"assignment_updated"`
}
