package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
)

var Database *gorm.DB

func Connect() (*gorm.DB, error) {
	var err error
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port, _ := strconv.Atoi(os.Getenv("DB_PORT"))

	createDBDsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/", username, password, host, port)
	database, err := gorm.Open(mysql.Open(createDBDsn), &gorm.Config{})

	_ = database.Exec("CREATE DATABASE IF NOT EXISTS " + databaseName + ";")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", username, password, host, port, databaseName)
	Database, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Print("Database connection failed: ", err)
		return database, err
	} else {
		log.Print("Successfully connected to the database")
		return database, nil
	}
}
