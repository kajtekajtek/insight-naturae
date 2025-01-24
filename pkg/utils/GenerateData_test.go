// pkg/utils/GenerateData_test.go - unit test for the GenerateData function
package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenerateData(t *testing.T) {
	min_max := []struct {
		min float64
		max float64
	}{
		{min: -10.0, max: 10.0},
		{min: -5.0, max: 5.0},
		{min: -100.0, max: 100.0},
		{min: -0.5, max: 0.5},
		{min: -1.0, max: 1.0},
		{min: -50.0, max: 50.0},
		{min: -1000.0, max: 1000.0},
		{min: -0.1, max: 0.1},
		{min: -0.01, max: 0.01},
		{min: -0.001, max: 0.001},
	}

	for _, mm := range min_max {
		data := GenerateData(mm.min, mm.max)
		assert.GreaterOrEqual(t, data, mm.min)
		assert.LessOrEqual(t, data, mm.max)
	}
}