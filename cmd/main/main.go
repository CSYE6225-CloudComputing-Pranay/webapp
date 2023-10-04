package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"webapp/controller"
	"webapp/database"
)

func main() {
	loadEnv()
	loadDatabaseScripts()
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		log.Print(".env files not found", err)
	} else {
		log.Print("Environment variables loaded successfully!!")
	}
}

func loadDatabaseScripts() {

	_, connectionError := database.Connect()
	if connectionError != nil {
		log.Fatal("Error connecting to the database", connectionError)
	}
	migrationError := database.Database.AutoMigrate(&database.Account{}, &database.Assignment{})
	if migrationError != nil {
		log.Print("Error while running auto migrate on the database", migrationError)
	}

	fileError := database.LoadDataFromFile(database.Database, os.Getenv("FILE_PATH"))
	if fileError != nil {
		log.Fatal("Error loading database scripts", fileError)
	}
	log.Print("Database scripts loaded successfully!!")
}

func serveApplication() {
	router := gin.Default()
	router.Use(DefaultHeaders())
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

	log.Print("Starting server with Gin framework")

	err := router.Run()
	if err != nil {
		log.Fatal("Error occurred while running starting server with Gin Framework", err)
	}
}
