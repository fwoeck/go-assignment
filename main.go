package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
)

func main() {
	db := initDB()

	router := gin.Default()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/matches/flooring", func(c *gin.Context) {
		Flooring(c, db)
	})

	router.GET("/partners", func(c *gin.Context) {
		var partners []Partners
		db.Preload("Services").Find(&partners)
		c.JSON(http.StatusOK, gin.H{"partners": partners})
	})

	router.Run(":8080")
}
