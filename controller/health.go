package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"webapp/database"
	"webapp/logger"
)

func Health(context *gin.Context) {

	logger.GetMetricsClient().Incr("assignment.health", 1)

	var writer = context.Writer

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		zap.L().Error("Request contains unwanted request body or query parameters")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlDB, err := database.Database.DB()
	if err != nil {
		zap.L().Error("Error while connecting to database", zap.Error(err))
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		zap.L().Error("Error while trying to ping database", zap.Error(err))
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}
