// pkg/utils/GenerateSecretKey.go - encryption secret key generating utility
package utils

import (
	"crypto/rand"
	"fmt"
)

// generate a length byte long random secret key
func GenerateSecretKey(length int) ([]byte, error) {
	if length < 0 {
		return nil, fmt.Errorf("Invalid key length: %d", length)
	}	

	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return nil, fmt.Errorf("Failed while generating secret key: %v", err)
	}

	return key, nil
}