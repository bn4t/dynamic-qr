package utils

import (
	"os"
	"path/filepath"
)

// get an environment variable and use the fallback if it is undefined
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// get the current execution dir
func GetExecutionDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(ex), nil
}
