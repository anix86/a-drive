package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupFileRoutes(router *gin.RouterGroup) {
	router.GET("/files", handlers.ListFiles)
	router.GET("/photos", handlers.GetPhotos)
	router.POST("/files/upload", handlers.UploadFile)
	router.GET("/files/:id/download", handlers.DownloadFile)
	router.DELETE("/files/:id", handlers.DeleteFile)
	router.PUT("/files/:id", handlers.RenameFile)
	
	// Search functionality
	router.GET("/search", handlers.SearchFiles)
	router.GET("/files/types", handlers.GetFileTypes)
	
	// Bulk operations
	router.POST("/bulk", handlers.BulkOperation)
	
	// File versioning
	router.POST("/files/:id/versioning/enable", handlers.EnableVersioning)
	router.POST("/files/:id/versioning/disable", handlers.DisableVersioning)
	router.GET("/files/:id/versions", handlers.GetFileVersions)
	router.POST("/files/:id/versions", handlers.CreateNewVersion)
	router.POST("/files/:id/versions/:version_id/restore", handlers.RestoreVersion)
	router.GET("/files/:id/versions/:version_id/download", handlers.DownloadVersion)
}