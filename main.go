package main

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
)

// Define a struct for the expected JSON object
type QueryParams struct {
	AddressLon   float64   `json:"address_lon" binding:"required"`
	AddressLat   float64   `json:"address_lat" binding:"required"`
	Services     []string  `json:"services" binding:"required"`
	FloorSize    float64   `json:"floor_size" binding:"required"`
	PhoneNumber  string    `json:"phone_number" binding:"required"`
}

// Validate if every service in the list is valid
func validateServices(services []string) bool {
	validServices := map[string]bool{"wood": true, "carpet": true, "tiles": true}
	for _, service := range services {
		if _, ok := validServices[service]; !ok {
			return false
		}
	}
	return true
}

func main() {
	// Initialize Gin router
	router := gin.Default()

	// Define the GET endpoint
	router.GET("/matches/flooring", func(c *gin.Context) {
		// Extract the 'q' query parameter
		q := c.Query("q")

		// Deserialize the JSON object
		var queryParams QueryParams
		if err := json.Unmarshal([]byte(q), &queryParams); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter"})
			return
		}

		// Manual validation for services
		if !validateServices(queryParams.Services) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service(s) provided"})
			return
		}

		// Log the valid query parameters to a file
		logToFile(queryParams)

		// Respond to the client
		c.JSON(http.StatusOK, gin.H{"message": "Query parameters received and logged"})
	})

	// Run the server
	router.Run(":8080")
}

// Function to log the query parameters to a file
func logToFile(data QueryParams) {
	logFile, err := os.OpenFile("queries.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Received query: %+v\n", data)
}
