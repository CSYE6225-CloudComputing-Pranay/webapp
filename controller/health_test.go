package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"log"
	"net/http/httptest"
	"testing"
)

type HealthTestSuite struct {
	suite.Suite
	App *gin.Engine
}

func TestHealthTestSuite(t *testing.T) {
	suite.Run(t, &HealthTestSuite{})
}

func (s *HealthTestSuite) SetupSuite() {

	loadEnv()

	app := gin.New()
	app.GET("/healthz", Health)
	s.App = app
}

func loadEnv() {
	err := godotenv.Load(".env.test")
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}
	log.Print("Environment variables loaded successfully!!")
}

func (s *HealthTestSuite) TestHealthIntegrationTest() {

	req := httptest.NewRequest("GET", "/healthz", nil)
	w := httptest.NewRecorder()
	s.App.ServeHTTP(w, req)

	res := w.Code

	s.Equal(200, res)
}
