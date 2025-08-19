package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupAdminRoutes(router *gin.RouterGroup) {
	router.GET("/users", handlers.ListUsers)
	router.POST("/users", handlers.CreateUser)
	router.GET("/files", handlers.BrowseUserFiles)
	
	// Configuration endpoints
	router.GET("/config", handlers.GetConfig)
	router.GET("/cors", handlers.CORSInfo)
}