package controller

import (
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
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
			if err := database.Database.Where("email = ?", username).First(&account).Error; err != nil {
				log.Print("Error while fetching the account, missing email: ", err)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
				log.Print("Error while fetching the account, missing password: ", err)
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Set("email", account.Email)
			return
		}

		log.Print("Error while fetching the account, missing token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	})
}
