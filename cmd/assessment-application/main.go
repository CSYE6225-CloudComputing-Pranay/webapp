package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"net/http"
	"os"
	"webapp/controller"
	"webapp/database"
	"webapp/logger"
)

func main() {
	loadEnv()
	log := logger.InitLogger()
	client := logger.InitMetrics()
	defer log.Sync()
	defer client.Close()
	setupDatabase()
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		zap.L().Error("error while fetching default environment variables: ", zap.Error(err))
	} else {
		zap.L().Info("Environment variables loaded successfully!!")
	}
}

func setupDatabase() {

	_, connectionError := database.Connect()
	if connectionError != nil {
		zap.L().Error("Error connecting to the database", zap.Error(connectionError))
	} else {
		migrationError := database.Database.AutoMigrate(&database.Account{}, &database.Assignment{}, &database.Submission{})
		if migrationError != nil {
			zap.L().Fatal("Error while running auto migrate on the database", zap.Error(migrationError))
		} else {
			fileError := database.LoadDataFromFile(database.Database, os.Getenv("FILE_PATH"))
			if fileError != nil {
				zap.L().Warn("Error loading database scripts", zap.Error(fileError))
			} else {
				zap.L().Info("Database scripts loaded successfully!!")
			}
		}
	}
}

func serveApplication() {
	router := gin.Default()
	router.Use(DefaultHeaders(), LogRequestResponse())

	binding.EnableDecoderDisallowUnknownFields = true
	router.HandleMethodNotAllowed = true

	router.NoMethod(func(c *gin.Context) {
		var writer = c.Writer
		writer.WriteHeader(http.StatusMethodNotAllowed)
		writer.Write(nil)
	})

	publicRoutes := router.Group("")
	publicRoutes.GET("/healthz", controller.Health)

	privateRoutes := router.Group("/v1/assignments")
	privateRoutes.Use(controller.BasicAuth())
	privateRoutes.POST("", controller.CreateAssignment)
	privateRoutes.GET("", controller.GetAllAssignments)
	privateRoutes.GET("/:assignmentID", controller.GetAssignment)
	privateRoutes.PUT("/:assignmentID", controller.UpdateAssignment)
	privateRoutes.DELETE("/:assignmentID", controller.DeleteAssignment)
	privateRoutes.POST("/:assignmentID/submission", controller.SubmitAssignment)
	privateRoutes.GET("/account", controller.CreateAccount)

	zap.L().Info("Starting server with Gin framework")

	err := router.Run()
	if err != nil {
		zap.L().Fatal("Error occurred while running starting server with Gin Framework", zap.Error(err))
	}

}
