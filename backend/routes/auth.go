package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

// Public auth routes (no authentication required)
func SetupAuthRoutes(router *gin.RouterGroup) {
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
}

// Protected auth routes (authentication required)
func SetupProtectedAuthRoutes(router *gin.RouterGroup) {
	router.GET("/auth/me", handlers.GetMe)
	router.GET("/profile", handlers.GetProfile)
	router.PUT("/profile", handlers.UpdateProfile)
	router.POST("/profile/change-password", handlers.ChangePassword)
}