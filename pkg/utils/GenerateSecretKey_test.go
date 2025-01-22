package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSecretKey(t *testing.T) {
	lengths := []int{-1, 0, 1, 16, 32, 64}

	for _, length := range lengths {
		key, err := GenerateSecretKey(length)
		if length < 0 {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
			assert.Len(t, key, length)
		}
	}
}