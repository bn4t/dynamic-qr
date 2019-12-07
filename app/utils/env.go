package utils

import "os"

// get an environment variable and use the fallback if it is undefined
func Getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
