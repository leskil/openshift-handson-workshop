package config

import (
	"os"
	"errors"
)

// AuthKey reads the environment variable AUTH_KEY or returns an error.
func AuthKey() (string, error) {
	key := os.Getenv("AUTH_KEY")

	if key != "" {
		return key, nil
	}

	return "", errors.New("Environment variable AUTH_KEY does not exist")
}