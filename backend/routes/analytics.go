package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupAnalyticsRoutes(router *gin.RouterGroup) {
	router.GET("/analytics/system", handlers.GetSystemAnalytics)
	router.GET("/analytics/user", handlers.GetUserAnalytics)
}