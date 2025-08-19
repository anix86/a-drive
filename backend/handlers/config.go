package handlers

import (
	"net/http"

	"a-drive-backend/config"

	"github.com/gin-gonic/gin"
)

// ConfigResponse represents the configuration data returned to admin
type ConfigResponse struct {
	CORSOrigins []string `json:"cors_origins"`
	CORSMethods []string `json:"cors_methods"`
	CORSHeaders []string `json:"cors_headers"`
	Port        string   `json:"port"`
	MaxFileSize int64    `json:"max_file_size"`
}

// GetConfig returns current configuration (admin only)
func GetConfig(c *gin.Context) {
	cfg := config.Load()

	response := ConfigResponse{
		CORSOrigins: cfg.GetCORSOrigins(),
		CORSMethods: cfg.GetCORSMethods(),
		CORSHeaders: cfg.GetCORSHeaders(),
		Port:        cfg.Port,
		MaxFileSize: cfg.MaxFileSize,
	}

	c.JSON(http.StatusOK, gin.H{
		"config": response,
		"message": "Configuration retrieved successfully",
	})
}

// CORSInfo provides information about CORS configuration
func CORSInfo(c *gin.Context) {
	cfg := config.Load()

	c.JSON(http.StatusOK, gin.H{
		"cors": gin.H{
			"origins": cfg.GetCORSOrigins(),
			"methods": cfg.GetCORSMethods(), 
			"headers": cfg.GetCORSHeaders(),
		},
		"message": "CORS configuration retrieved successfully",
		"note":    "To modify CORS settings, update the .env file and restart the server",
	})
}