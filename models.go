package main

import (
	"gorm.io/gorm"
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
	// For local calculation with the haversine function:
	Distance         float64     `gorm:"-"`
}

type QueryParams struct {
	AddressLon   float64   `json:"address_lon" binding:"required"`
	AddressLat   float64   `json:"address_lat" binding:"required"`
	Services     []string  `json:"services" binding:"required"`
	FloorSize    float64   `json:"floor_size" binding:"required"`
	PhoneNumber  string    `json:"phone_number" binding:"required"`
}
