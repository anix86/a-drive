package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabasePath    string
	JWTSecret      string
	RootDirectory  string
	MaxFileSize    int64
	AllowedTypes   string
	Port          string
}

func Load() *Config {
	maxSize, _ := strconv.ParseInt(getEnv("MAX_FILE_SIZE", "104857600"), 10, 64)
	
	return &Config{
		DatabasePath:   getEnv("DATABASE_PATH", "./storage/database.db"),
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key"),
		RootDirectory: getEnv("ROOT_DIRECTORY", "./storage/files"),
		MaxFileSize:   maxSize,
		AllowedTypes:  getEnv("ALLOWED_FILE_TYPES", "*"),
		Port:         getEnv("PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}