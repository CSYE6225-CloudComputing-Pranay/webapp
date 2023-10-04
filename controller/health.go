package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"webapp/database"
)

func Health(context *gin.Context) {

	var writer = context.Writer

	if context.Request.Body != http.NoBody || len(context.Request.URL.Query()) != 0 {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	if _, err := database.Connect(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	writer.WriteHeader(http.StatusOK)
	return
}
