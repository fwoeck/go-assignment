package main

import (
	"encoding/json"
	"net/http"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
)

// @title Flooring Matches API
// @description API for retrieving flooring matches based on user preferences.
// @version 1
// @host localhost:8080
// @BasePath /
// @Summary Retrieve flooring matches
// @Description Retrieves matches based on flooring preferences
// @Tags matches
// @Accept json
// @Produce json
// @Param q query string true "QueryParams Object" example('{"address_lon": 13.45, "address_lat": 52.5, "services": ["tiles", "carpet"], "floor_size": 120.5, "phone_number": "123-456-7890"}')
// @Success 200 {object} QueryParams
// @Failure 400 {object} map[string]string
// @Router /matches/flooring [get]
func Flooring(c *gin.Context, db *gorm.DB) {
	var queryParams QueryParams

	q := c.Query("q")

	if err := json.Unmarshal([]byte(q), &queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameter"})
		return
	}

	if !validateServices(queryParams.Services) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid service(s) provided"})
		return
	}

	logToFile(queryParams);

	partners := fetchPartners(queryParams, db);
	c.JSON(http.StatusOK, gin.H{"Partners": partners})
}

func PartnerIndex(c *gin.Context, db *gorm.DB) {
	var partners []Partners

	db.Preload("Services").Find(&partners)
	c.JSON(http.StatusOK, gin.H{"Partners": partners})
}
