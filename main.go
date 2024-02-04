package main

import (
	"encoding/json"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
	_ "matching/docs"
)

type Services struct {
	ID         uint    `gorm:"primaryKey"`
	PartnerID  uint
	Name       string  `gorm:"type:varchar(100)"`
}

type Partners struct {
	gorm.Model
	AddressLon       float64     `gorm:"type:float"`
	AddressLat       float64     `gorm:"type:float"`
	OperatingRadius  float64     `gorm:"type:float"`
	Rating           float64     `gorm:"type:float"`
	Services         []Services  `gorm:"foreignKey:PartnerID"`
	Distance         float64     `gorm:"-"`
}

type QueryParams struct {
	AddressLon   float64   `json:"address_lon" binding:"required"`
	AddressLat   float64   `json:"address_lat" binding:"required"`
	Services     []string  `json:"services" binding:"required"`
	FloorSize    float64   `json:"floor_size" binding:"required"`
	PhoneNumber  string    `json:"phone_number" binding:"required"`
}

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("matching.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	db.AutoMigrate(&Partners{}, &Services{})

	return db
}

func validateServices(services []string) bool {
	validServices := map[string]bool{"wood": true, "carpet": true, "tiles": true}
	for _, service := range services {
		if _, ok := validServices[service]; !ok {
			return false
		}
	}
	return true
}

func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	var (
		rad      = math.Pi / 180
		r        = 6378100.0
		dLat     = (lat2 - lat1) * rad
		dLon     = (lon2 - lon1) * rad
		a        = math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(lat1*rad)*math.Cos(lat2*rad)*math.Sin(dLon/2)*math.Sin(dLon/2)
		c        = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
		distance = r * c
	)

	return distance
}

func filterAndSortPartners(partners []Partners, queryParams QueryParams) []Partners {
	filteredPartners := make([]Partners, 0)

	for _, partner := range partners {
		distance := haversine(queryParams.AddressLat, queryParams.AddressLon, partner.AddressLat, partner.AddressLon)
		if distance <= float64(partner.OperatingRadius) {
			partner.Distance = distance
			filteredPartners = append(filteredPartners, partner)
		}
	}

	sort.Slice(filteredPartners, func(i, j int) bool {
		if filteredPartners[i].Rating == filteredPartners[j].Rating {
			return filteredPartners[i].Distance > filteredPartners[j].Distance
		}
		return filteredPartners[i].Rating > filteredPartners[j].Rating
	})

	return filteredPartners
}

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
// @Param q query string true "QueryParams Object" example("{\"address_lon\": 10.0, \"address_lat\": 20.0, \"services\": [\"wood\", \"tiles\"], \"floor_size\": 120.5, \"phone_number\": \"123-456-7890\"}")
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

	var partnerIDs []uint

	db.Table("services").
		Select("services.partner_id").
		Where("services.name IN ?", queryParams.Services).
		Group("services.partner_id").
		Having("COUNT(DISTINCT services.name) >= ?", len(queryParams.Services)).
		Pluck("services.partner_id", &partnerIDs)

	var partners []Partners

	if len(partnerIDs) > 0 {
		db.Preload("Services").
			Where("id IN ?", partnerIDs).
			Order("rating DESC").
			Limit(1000).
			Find(&partners)
	}

	partners = filterAndSortPartners(partners, queryParams)
	c.JSON(http.StatusOK, partners)
}

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

func logToFile(data QueryParams) {
	logFile, err := os.OpenFile("queries.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Received query: %+v\n", data)
}
