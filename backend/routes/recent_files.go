package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupRecentFilesRoutes(r *gin.RouterGroup) {
	r.GET("/recent-files", handlers.GetRecentFiles)
	r.POST("/recent-files/track/file/:id", handlers.TrackFileAccess)
	r.POST("/recent-files/track/folder/:id", handlers.TrackFolderAccess)
}