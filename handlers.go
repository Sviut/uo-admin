package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func addDeliveryHandler(c *gin.Context) {
	var req DeliveryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deliveryTime, err := time.Parse(time.RFC3339, req.Time)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid time format"})
		return
	}

	for _, oreReq := range req.Ores {
		delivery := OreDelivery{
			Color:             oreReq.Color,
			Quantity:          oreReq.Quantity,
			DeliveryTimestamp: deliveryTime,
		}

		fmt.Println(delivery)

		if err := db.Create(&delivery).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add delivery information"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"success": "Delivery information added successfully"})
}

func getDeliveriesHandler(c *gin.Context) {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)

	var deliveries []OreDelivery
	result := db.Where("delivery_timestamp >= ?", sevenDaysAgo).Find(&deliveries)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query delivery information"})
		return
	}

	response := make([]map[string]interface{}, 0)
	for _, delivery := range deliveries {
		dateTimeStr := delivery.DeliveryTimestamp.Format("1999")
		response = append(response, map[string]interface{}{
			"dateTime": dateTimeStr,
			"ores": []map[string]int{
				{delivery.Color: delivery.Quantity},
			},
		})
	}

	c.JSON(http.StatusOK, response)
}
