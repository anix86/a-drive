package routes

import (
	"github.com/gin-gonic/gin"
	"a-drive-backend/handlers"
)

func SetupFavoriteRoutes(r *gin.RouterGroup) {
	r.GET("/favorites", handlers.GetFavorites)
	r.POST("/favorites", handlers.AddFavorite)
	r.DELETE("/favorites/:id", handlers.RemoveFavorite)
	r.DELETE("/favorites/item", handlers.RemoveFavoriteByItem)
	r.GET("/favorites/check/:type/:id", handlers.CheckFavorite)
}