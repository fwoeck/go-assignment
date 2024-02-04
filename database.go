package main

import (
	"log"
	"sort"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func initDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("matching.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database", err)
	}

	db.AutoMigrate(&Partners{}, &Services{})

	return db
}

// This filters and sorts the partners objects in memory.
// See README.md for considerations.
//
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

// Retrieve only partner entries that have at least the requested
// services [tiles, carpet, wood] associated.
// For practical reasons, the response is limited to protect against
// overly large results. See README.md for considerations.
//
func fetchPartners(queryParams QueryParams, db *gorm.DB) []Partners {
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

	return filterAndSortPartners(partners, queryParams)
}
