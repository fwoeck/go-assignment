package main

import (
	"log"
	"math"
	"os"
)

// Standard function to calculate the distance of two points
// given as lon/lat. See README.md for considerations.
//
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

func logToFile(data QueryParams) {
	logFile, err := os.OpenFile("queries.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer logFile.Close()

	logger := log.New(logFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Printf("Received query: %+v\n", data)
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
