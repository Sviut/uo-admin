package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func loadEnvVariables() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}

func connectDB() {
	dbHost := os.Getenv("DB_HOST")
	pgUser := os.Getenv("POSTGRES_USER")
	pgPass := os.Getenv("POSTGRES_PASSWORD")

	if dbHost == "" || pgUser == "" || pgPass == "" {
		log.Fatal("One or more required environment variables are not set")
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=5432 sslmode=disable", dbHost, pgUser, pgPass)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
}

func migrateDB() {
	if err := db.AutoMigrate(&OreDelivery{}); err != nil {
		log.Fatal("Failed to perform auto migration:", err)
	}
}

func initDB() {
	loadEnvVariables()
	connectDB()
	migrateDB()
}

func main() {
	initDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World test!")
	})

	r.POST("/deliveries", addDeliveryHandler)
	r.GET("/deliveries", getDeliveriesHandler)

	r.Run(":8080")
}
