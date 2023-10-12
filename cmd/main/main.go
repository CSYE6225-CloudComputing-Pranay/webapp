package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"webapp/controller"
	"webapp/database"
)

func main() {
	loadEnv()
	setupDatabase()
	serveApplication()
}

func loadEnv() {
	err := godotenv.Load(".env.local", ".env")
	if err != nil {
		log.Print(".env files not found: ", err)
	} else {
		log.Print("Environment variables loaded successfully!!")
	}
}

func setupDatabase() {

	_, connectionError := database.Connect()
	if connectionError != nil {
		log.Println("Error connecting to the database: ", connectionError)
	} else {
		migrationError := database.Database.AutoMigrate(&database.Account{}, &database.Assignment{})
		if migrationError != nil {
			log.Fatal("Error while running auto migrate on the database: ", migrationError)
		} else {
			fileError := database.LoadDataFromFile(database.Database, os.Getenv("FILE_PATH"))
			if fileError != nil {
				log.Println("Error loading database scripts: ", fileError)
			} else {
				log.Print("Database scripts loaded successfully!!")
			}
		}
	}
}

func serveApplication() {
	router := gin.Default()
	router.Use(DefaultHeaders())

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

	log.Print("Starting server with Gin framework")

	err := router.Run()
	if err != nil {
		log.Fatal("Error occurred while running starting server with Gin Framework", err)
	}

}
