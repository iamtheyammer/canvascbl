package services

import "math"

const postgresMaxPlaceholders = 65535

// CalculateChunkSize calculates the size of each insert chunk to fit under postgresMaxPlaceholders.
func CalculateChunkSize(placeholdersPerValue int) int {
	return int(math.Floor(float64(postgresMaxPlaceholders / placeholdersPerValue)))
}
