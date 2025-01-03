// pkg/utils/GenerateData.go - float number generating utility
package utils

import (
	"math/rand"
)

// GenerateData generates a random float64 number between min and max
func GenerateData(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}