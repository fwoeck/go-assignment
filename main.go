package main

import (
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
		PartnerIndex(c, db)
	})

	router.Run(":8080")
}
