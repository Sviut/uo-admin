package main

import (
	"fmt"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"net/http"
	"sort"
	"strings"
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
	result := db.
		Not("color = ?", "Iron Ore").
		Where("delivery_timestamp >= ?", sevenDaysAgo).
		Find(&deliveries)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query delivery information"})
		return
	}

	timestampSet := make(map[string]bool)
	for _, d := range deliveries {
		timestamp := d.DeliveryTimestamp.Format("02 15:04")
		timestampSet[timestamp] = true
	}

	var timestamps []string
	for timestamp := range timestampSet {
		timestamps = append(timestamps, timestamp)
	}
	sort.Strings(timestamps)

	// Подготовка данных для графика
	colorDataMap := make(map[string][]opts.LineData)
	for _, d := range deliveries {
		if _, ok := colorDataMap[d.Color]; !ok {
			colorDataMap[d.Color] = make([]opts.LineData, len(timestamps))
			for i := range colorDataMap[d.Color] {
				colorDataMap[d.Color][i] = opts.LineData{Value: 0} // Инициализация нулями
			}
		}
		timestampIndex := sort.SearchStrings(timestamps, d.DeliveryTimestamp.Format("02 15:04"))
		colorDataMap[d.Color][timestampIndex] = opts.LineData{Value: d.Quantity}
	}

	// Создание графика 1
	lineChart := charts.NewLine()
	lineChart.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true, Height: "400"}),
		charts.WithYAxisOpts(opts.YAxis{Type: "value"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true, Trigger: "axis"}),
	)
	lineChart.SetXAxis(timestamps)

	for color, data := range colorDataMap {
		lineChart.AddSeries(color, data)
	}

	css := `<style>
                .charts {
                    margin-top: 160px;
                }
            </style>`

	var b strings.Builder
	lineChart.Render(&b)

	c.Header("Content-Type", "text/html")
	c.Writer.Write([]byte(css + "<div class='charts'>" + b.String() + "</div>"))
}
