// pkg/utils/Getenv_test.go - unit test for the Getenv function
package utils

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
)

func TestGetenv(t *testing.T) {
	key := "TEST_KEY"
	fallback := "fallback"

	t.Run("when the environment variable is set", func(t *testing.T) {
		value := "value"
		os.Setenv(key, value)
		assert.Equal(t, value, Getenv(key, fallback))
	})

	t.Run("when the environment variable is not set", func(t *testing.T) {
		os.Unsetenv(key)
		assert.Equal(t, fallback, Getenv(key, fallback))
	})
}
