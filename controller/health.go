package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/database"
	"webapp/logger"
)

func Health(context *gin.Context) {

	logger.GetMetricsClient().Incr("assignment.health", 1)

	var writer = context.Writer

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	sqlDB, err := database.Database.DB()
	if err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}
