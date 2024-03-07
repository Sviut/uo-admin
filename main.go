package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"os"
)

var db *gorm.DB

func initDB() {
	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dsn := fmt.Sprintf("host=%s user=root password=root dbname=postgres port=5432 sslmode=disable", dbHost)

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := db.AutoMigrate(&OreDelivery{}); err != nil {
		log.Fatal("Failed to perform auto migration:", err)
	}
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
