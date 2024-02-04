package main

import (
	"testing"
	"gorm.io/gorm"
)

func TestFilterAndSortPartners(t *testing.T) {
	queryParams := QueryParams{
		AddressLon: 13.4050,
		AddressLat: 52.5200,
	}

	partners := []Partners{
		// Within radius:
		{Model: gorm.Model{ID: 1}, AddressLon: 13.4050, AddressLat: 52.5200, OperatingRadius: 5000, Rating: 4.5, Services: []Services{}},
		{Model: gorm.Model{ID: 2}, AddressLon: 13.4550, AddressLat: 52.5200, OperatingRadius: 10000, Rating: 4.7, Services: []Services{}},
		// Outside of radius:
		{Model: gorm.Model{ID: 3}, AddressLon: 14.4050, AddressLat: 53.5200, OperatingRadius: 1000, Rating: 4.9, Services: []Services{}},
	}

	expectedPartnerIDs := []uint{2, 1}

	filteredSortedPartners := filterAndSortPartners(partners, queryParams)

	if len(filteredSortedPartners) != len(expectedPartnerIDs) {
		t.Errorf("Expected %d partners, got %d", len(expectedPartnerIDs), len(filteredSortedPartners))
	}

	for i, partner := range filteredSortedPartners {
		if partner.ID != expectedPartnerIDs[i] {
			t.Errorf("Expected partner ID %d at position %d, got %d", expectedPartnerIDs[i], i, partner.ID)
		}
	}
}
