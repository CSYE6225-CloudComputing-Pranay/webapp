package controller

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm/clause"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"webapp/database"
	"webapp/logger"
)

func CreateAssignment(context *gin.Context) {

	var request AssignmentRequest

	logger.GetMetricsClient().Incr("assignment.create.record", 1)

	if len(context.Request.URL.Query()) != 0 {
		zap.L().Error("Request contains unwanted request query parameters", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery))
		context.Status(http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().UTC().After(request.Deadline) {
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

	logger.GetMetricsClient().Incr("assignment.fetch.records", 1)

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		requestBody, err := io.ReadAll(context.Request.Body)
		defer context.Request.Body.Close()
		if err != nil {
			zap.L().Error("Error while closing the request body", zap.Error(err))
			requestBody = []byte("Malformed request body")
		}

		zap.L().Error("Request contains unwanted request query parameters or body", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery), zap.ByteString("request-body", requestBody))
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

	logger.GetMetricsClient().Incr("assignment.fetch.record", 1)

	var assignment database.Assignment

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		requestBody, err := io.ReadAll(context.Request.Body)
		defer context.Request.Body.Close()
		if err != nil {
			zap.L().Error("Error while closing the request body", zap.Error(err))
			requestBody = []byte("Malformed request body")
		}
		zap.L().Error("Request contains unwanted request query parameters or body", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery), zap.ByteString("request-body", requestBody))
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

	logger.GetMetricsClient().Incr("assignment.update.record", 1)

	var request AssignmentRequest

	if len(context.Request.URL.Query()) != 0 {
		zap.L().Error("Request contains unwanted request query parameters", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery))
		context.Status(http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().UTC().After(request.Deadline) {
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

	logger.GetMetricsClient().Incr("assignment.delete.record", 1)

	var assignment database.Assignment

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		requestBody, err := io.ReadAll(context.Request.Body)
		defer context.Request.Body.Close()
		if err != nil {
			requestBody = []byte("Malformed request body")
			zap.L().Error("Error while closing the request body", zap.Error(err))
		}
		zap.L().Error("Request contains unwanted request query parameters or body", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery), zap.ByteString("request-body", requestBody))
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

	database.Database.Select(clause.Associations).Delete(&assignment)
	context.Status(http.StatusNoContent)
	return
}

func SubmitAssignment(context *gin.Context) {

	var request SubmitRequest
	assignmentID := context.Param("assignmentID")
	accountID := context.GetString("accountID")

	logger.GetMetricsClient().Incr("assignment.submit.record", 1)

	if len(context.Request.URL.Query()) != 0 {
		zap.L().Error("Request contains unwanted request query parameters", zap.String("user-mail", context.GetString("email")),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.String("request-query", context.Request.URL.RawQuery))
		context.Status(http.StatusBadRequest)
		return
	}

	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var count int64
	if err := database.Database.Model(&database.Submission{}).Where(database.Submission{AssignmentID: assignmentID, AccountID: accountID}).
		Count(&count).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		zap.L().Error("Failed to get submission count from database", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path), zap.Error(err))
		return
	}

	var account database.Account
	account.ID = accountID

	if err := database.Database.First(&account).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}

	var assignment database.Assignment
	if err := database.Database.Where(database.Assignment{ID: assignmentID}).First(&assignment).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		zap.L().Error("Failed to get assignment from database for number of attempts", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path))
	}

	if time.Now().UTC().After(assignment.Deadline) {
		zap.L().Error("User attempted submission after the deadline", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path))
		topicARN, msg := prepareSubmissionMessage("DEADLINE_PASSED", account, request, assignmentID, count+1)
		publishMessage(topicARN, msg)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Deadline has passed"})
		return
	}

	if count >= int64(assignment.NumOfAttempts) {
		zap.L().Error("User reached maximum number of attempts", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path))
		topicARN, msg := prepareSubmissionMessage("MAX_ATTEMPTS", account, request, assignmentID, count+1)
		publishMessage(topicARN, msg)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Maximum number of attempts reached"})
		return
	}

	zipSize, isValid := isValidZIP(request.SubmissionURL)

	if !isValid {
		zap.L().Error("User submitted invalid zip file", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path))
		topicARN, msg := prepareSubmissionMessage("INVALID_URL", account, request, assignmentID, count+1)
		publishMessage(topicARN, msg)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid zip file"})
		return
	}

	if zipSize == "0" {
		zap.L().Error("User submitted empty zip file", zap.String("user-mail", context.GetString("email")), zap.String("assignment-id", assignmentID), zap.String("account-id", accountID),
			zap.String("request-method", context.Request.Method), zap.String("request-path", context.Request.URL.Path))
		topicARN, msg := prepareSubmissionMessage("NO_CONTENT", account, request, assignmentID, count+1)
		publishMessage(topicARN, msg)
		context.JSON(http.StatusBadRequest, gin.H{"error": "Empty zip file"})
		return
	}

	id := uuid.New().String()

	submission := database.Submission{
		ID:            id,
		SubmissionURL: request.SubmissionURL,
		AssignmentID:  assignmentID,
		AccountID:     accountID,
	}

	if err := database.Database.Create(&submission).Error; err != nil {
		context.Status(http.StatusServiceUnavailable)
		return
	}

	topicARN, msg := prepareSubmissionMessage("SUCCESS", account, request, assignmentID, count+1)

	publishMessage(topicARN, msg)

	submissionResponse := SubmissionResponse{
		ID:                submission.ID,
		AssignmentID:      submission.AssignmentID,
		SubmissionURL:     submission.SubmissionURL,
		SubmissionDate:    submission.SubmissionDate,
		SubmissionUpdated: submission.SubmissionUpdated,
	}

	context.JSON(http.StatusCreated, submissionResponse)
	return
}

func prepareSubmissionMessage(status string, account database.Account, request SubmitRequest, assignmentID string, count int64) (string, string) {
	topicARN := os.Getenv("SUBMISSION_TOPIC_ARN")
	submissionMessage := SubmissionMessage{
		Status:        status,
		UserEmail:     account.Email,
		SubmissionURL: request.SubmissionURL,
		AssignmentID:  assignmentID,
		FirstName:     account.FirstName,
		LastName:      account.LastName,
		Count:         count}

	msgBytes, err := json.Marshal(submissionMessage)
	if err != nil {
		zap.L().Error("Error marshaling SNS message", zap.Error(err))
	}
	msg := string(msgBytes)
	return topicARN, msg
}

func publishMessage(topicARN string, msg string) {

	if msg == "" || topicARN == "" {
		zap.L().Error("Missing required fields in SNS message", zap.String("message", msg), zap.String("topic-arn", topicARN))
		return
	}

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		zap.L().Error("Error loading default AWS config", zap.Error(err))
		return
	}

	client := sns.NewFromConfig(cfg)

	input := &sns.PublishInput{
		Message:  &msg,
		TopicArn: &topicARN,
	}

	_, err = client.Publish(context.TODO(), input)
	if err != nil {
		zap.L().Error("Error publishing message to SNS topic", zap.Error(err))
		return
	}
}

func isValidZIP(url string) (string, bool) {
	resp, err := http.Head(url)
	if err != nil {
		return "0", false
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	contentLength := resp.Header.Get("Content-Length")
	return contentLength, strings.EqualFold(contentType, "application/zip")
}

type SubmissionMessage struct {
	Status        string `json:"status"`
	SubmissionURL string `json:"submissionUrl"`
	UserEmail     string `json:"userEmail"`
	AssignmentID  string `json:"assignmentId"`
	AccountID     string `json:"account_id"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Count         int64  `json:"attempt"`
}
