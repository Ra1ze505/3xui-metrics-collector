package main

import (
	"os"
	"strconv"
)

type Config struct {
	XUIHost     string
	XUIPort     string
	XUIBasePath string
	XUIUsername string
	XUIPassword string
	XUIUseTLS   bool
}

func LoadConfig() (*Config, error) {
	config := &Config{
		XUIHost:     getEnv("X_UI_HOST", ""),
		XUIPort:     getEnv("X_UI_PORT", ""),
		XUIBasePath: getEnv("X_UI_BASEPATH", ""),
		XUIUsername: getEnv("X_UI_USERNAME", ""),
		XUIPassword: getEnv("X_UI_PASSWORD", ""),
		XUIUseTLS:   getBoolEnv("X_UI_USE_TLS", true),
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

func getBoolEnv(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return boolValue
}
