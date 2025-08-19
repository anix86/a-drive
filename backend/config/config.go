package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DatabasePath    string
	JWTSecret      string
	RootDirectory  string
	MaxFileSize    int64
	AllowedTypes   string
	Port          string
	CORSOrigins    string
	CORSMethods    string
	CORSHeaders    string
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
		CORSOrigins:   getEnv("CORS_ORIGINS", "http://localhost:3000,http://localhost:3001"),
		CORSMethods:   getEnv("CORS_METHODS", "GET,POST,PUT,DELETE,OPTIONS"),
		CORSHeaders:   getEnv("CORS_HEADERS", "Origin,Content-Type,Authorization"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ParseCSV splits a comma-separated string into a slice
func (c *Config) ParseCSV(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	for i, part := range parts {
		parts[i] = strings.TrimSpace(part)
	}
	return parts
}

// GetCORSOrigins returns the CORS origins as a slice
func (c *Config) GetCORSOrigins() []string {
	return c.ParseCSV(c.CORSOrigins)
}

// GetCORSMethods returns the CORS methods as a slice  
func (c *Config) GetCORSMethods() []string {
	return c.ParseCSV(c.CORSMethods)
}

// GetCORSHeaders returns the CORS headers as a slice
func (c *Config) GetCORSHeaders() []string {
	return c.ParseCSV(c.CORSHeaders)
}