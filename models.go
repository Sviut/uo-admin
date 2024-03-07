package main

import "time"

type DeliveryRequest struct {
	Time string `json:"time"`
	Ores []struct {
		Color    string `json:"color"`
		Quantity int    `json:"quantity"`
	} `json:"ores"`
}

type OreDelivery struct {
	DeliveryID        uint   `gorm:"primaryKey"`
	Color             string `gorm:"not null"`
	Quantity          int    `gorm:"check:Quantity > 0"`
	DeliveryTimestamp time.Time
}
