package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupFolderRoutes(router *gin.RouterGroup) {
	router.POST("/folders", handlers.CreateFolder)
	router.GET("/folders/:id", handlers.GetFolder)
	router.GET("/folders/:id/breadcrumbs", handlers.GetFolderBreadcrumbs)
	router.PUT("/folders/:id", handlers.UpdateFolder)
	router.DELETE("/folders/:id", handlers.DeleteFolder)
	router.POST("/folders/:id/zip", handlers.CreateZipArchive)
}