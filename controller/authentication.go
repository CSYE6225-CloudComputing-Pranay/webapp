package controller

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"webapp/database"
)

func BasicAuth() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		writer := c.Writer
		writer.Header().Set("Content-Type", `application/json`)
		username, password, ok := c.Request.BasicAuth()
		if ok {

			var account database.Account

			sqlDB, err := database.Database.DB()
			if err != nil {
				zap.L().Error("Error while connecting to database", zap.Error(err))
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}

			if err := sqlDB.Ping(); err != nil {
				zap.L().Error("Error while connecting to database", zap.Error(err))
				c.AbortWithStatus(http.StatusServiceUnavailable)
				return
			}

			if err := database.Database.Where("email = ?", username).First(&account).Error; err != nil {
				zap.L().Error("Error while fetching the account, missing email", zap.Error(err))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
				zap.L().Error("Error while fetching the account, missing password", zap.Error(err))
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			zap.L().Info("User successfully authenticated", zap.String("user-mail", account.Email),
				zap.String("userId", account.ID), zap.String("method", c.Request.Method), zap.String("path", c.Request.URL.Path))

			c.Set("email", account.Email)
			c.Set("accountID", account.ID)
			return
		}

		zap.L().Error("Error while fetching the account, missing token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	})
}
