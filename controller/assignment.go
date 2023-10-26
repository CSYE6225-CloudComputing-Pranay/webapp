package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
	"webapp/database"
)

func CreateAssignment(context *gin.Context) {

	var request AssignmentRequest

	if len(context.Request.URL.Query()) != 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().After(request.Deadline) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Deadline is a past date"})
		return
	}

	id := uuid.New().String()

	assignment := database.Assignment{
		ID:            id,
		Name:          request.Name,
		Points:        request.Points,
		NumOfAttempts: request.NumOfAttempts,
		Deadline:      request.Deadline,
		AccountEmail:  context.GetString("email"),
	}

	if err := database.Database.Create(&assignment).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}

	assignmentResponse := AssignmentResponse{
		ID:                assignment.ID,
		Name:              assignment.Name,
		Points:            assignment.Points,
		NumOfAttempts:     assignment.NumOfAttempts,
		Deadline:          assignment.Deadline,
		AssignmentCreated: assignment.AssignmentCreated,
		AssignmentUpdated: assignment.AssignmentUpdated,
	}

	context.JSON(http.StatusCreated, assignmentResponse)
	return
}

func GetAllAssignments(context *gin.Context) {

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	var assignments []database.Assignment
	var assignmentResponses = make([]AssignmentResponse, 0)

	database.Database.Find(&assignments)

	for i := 0; i < len(assignments); i++ {
		assignmentResponses = append(assignmentResponses, AssignmentResponse{
			ID:                assignments[i].ID,
			Name:              assignments[i].Name,
			Points:            assignments[i].Points,
			NumOfAttempts:     assignments[i].NumOfAttempts,
			Deadline:          assignments[i].Deadline,
			AssignmentCreated: assignments[i].AssignmentCreated,
			AssignmentUpdated: assignments[i].AssignmentUpdated,
		})
	}

	context.JSON(http.StatusOK, assignmentResponses)
	return
}

func GetAssignment(context *gin.Context) {

	var assignment database.Assignment

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	id := context.Param("assignmentID")

	if err := database.Database.Where("id=?", id).First(&assignment).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	assignmentResponse := AssignmentResponse{
		ID:                assignment.ID,
		Name:              assignment.Name,
		Points:            assignment.Points,
		NumOfAttempts:     assignment.NumOfAttempts,
		Deadline:          assignment.Deadline,
		AssignmentCreated: assignment.AssignmentCreated,
		AssignmentUpdated: assignment.AssignmentUpdated,
	}

	context.JSON(http.StatusOK, assignmentResponse)
	return
}

func UpdateAssignment(context *gin.Context) {

	var request AssignmentRequest

	if len(context.Request.URL.Query()) != 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().After(request.Deadline) {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Deadline is a past date"})
		return
	}

	id := context.Param("assignmentID")
	email := context.GetString("email")

	var record database.Assignment

	assignment := database.Assignment{
		ID:            id,
		Name:          request.Name,
		Points:        request.Points,
		NumOfAttempts: request.NumOfAttempts,
		Deadline:      request.Deadline,
	}

	if err := database.Database.Where("id=?", id).First(&record).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	if record.AccountEmail != email {
		context.Status(http.StatusForbidden)
		return
	}

	if err := database.Database.Where(database.Assignment{ID: id}).Updates(&assignment).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		log.Print("error: ", err.Error())
		return
	}

	context.Status(http.StatusNoContent)
	return
}

func DeleteAssignment(context *gin.Context) {

	var assignment database.Assignment

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		context.Status(http.StatusBadRequest)
		return
	}

	id := context.Param("assignmentID")
	email := context.GetString("email")

	if err := database.Database.First(&assignment, database.Assignment{ID: id}).Error; err != nil {
		context.Status(http.StatusNotFound)
		return
	}

	if assignment.AccountEmail != email {
		context.Status(http.StatusForbidden)
		return
	}

	database.Database.Delete(&assignment)
	context.Status(http.StatusNoContent)
	return
}
