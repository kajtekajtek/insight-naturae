// pkg/utils/Getenv.go - utility function for getting environment variables
package utils

import (
	"os"
)

/* Getenv returns the value of the enviroment variable with the given name 
	or the fallback value if the variable is not set */
func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}