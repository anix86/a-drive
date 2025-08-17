package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupSharingRoutes(router *gin.RouterGroup) {
	// Protected sharing routes (require authentication)
	router.POST("/files/:id/share", handlers.CreateFileShare)
	router.GET("/files/:id/shares", handlers.GetFileShares)
	router.GET("/shares", handlers.GetUserShares)
	router.DELETE("/shares/:share_id", handlers.DeleteFileShare)
}

func SetupPublicSharingRoutes(router *gin.Engine) {
	// Public sharing routes (no authentication required)
	shareGroup := router.Group("/share")
	{
		shareGroup.POST("/:token/access", handlers.AccessSharedFile)
		shareGroup.GET("/:token/download", handlers.DownloadSharedFile)
		shareGroup.GET("/:token", handlers.AccessSharedFile) // For GET requests without password
	}
}