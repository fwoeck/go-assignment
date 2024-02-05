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
	Services     []string  `json:"services" binding:"required,gt=0"`
	FloorSize    float64   `json:"floor_size" binding:"required,gt=0"`
	PhoneNumber  string    `json:"phone_number" binding:"required"`
}

type FlooringMatchesResponse struct {
	Partners []Partner `json:"partners"`
}

type Partner struct {
	ID               uint       `json:"id"`
	AddressLon       float64    `json:"addressLon"`
	AddressLat       float64    `json:"addressLat"`
	OperatingRadius  float64    `json:"operatingRadius"`
	Rating           float64    `json:"rating"`
	Services         []Service  `json:"services"`
	Distance         float64    `json:"distance,omitempty"`
}

type Service struct {
	ID        uint   `json:"id"`
	Name      string `json:"name"`
}
