package main

import (
	"net/http"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	_ "matching/docs"
)

// PostFlooringMatches godoc
// @Summary Retrieve flooring matches
// @Description Retrieves matches based on flooring preferences submitted by the customer
// @Tags matches
// @Accept json
// @Produce json
// @Param query body QueryParams true "Query Parameters"
// @Success 200 {object} FlooringMatchesResponse "success"
// @Failure 400 {object} map[string]string "bad request, invalid input"
// @Router /matches/flooring [post]
func Flooring(c *gin.Context, db *gorm.DB) {
	var queryParams QueryParams
	if err := c.ShouldBindJSON(&queryParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
