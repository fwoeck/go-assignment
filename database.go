package main

import (
	"log"
	"sort"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "matching/docs"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("matching.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	db.AutoMigrate(&Partners{}, &Services{})

	return db
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
