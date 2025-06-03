package main

import (
	"os"
)

type Config struct {
	XUIHost     string
	XUIPort     string
	XUIBasePath string
	XUIUsername string
	XUIPassword string
}

func LoadConfig() (*Config, error) {
	config := &Config{
		XUIHost:     getEnv("X_UI_HOST", ""),
		XUIPort:     getEnv("X_UI_PORT", ""),
		XUIBasePath: getEnv("X_UI_BASEPATH", ""),
		XUIUsername: getEnv("X_UI_USERNAME", ""),
		XUIPassword: getEnv("X_UI_PASSWORD", ""),
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
