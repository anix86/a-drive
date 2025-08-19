package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"a-drive-backend/config"
	"a-drive-backend/database"
	"a-drive-backend/middleware"
	"a-drive-backend/routes"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.Load()
	db := database.Init(cfg.DatabasePath)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.GetCORSOrigins(),
		AllowMethods:     cfg.GetCORSMethods(),
		AllowHeaders:     cfg.GetCORSHeaders(),
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.DatabaseMiddleware(db))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "A-Drive is running"})
	})

	// CORS info endpoint (public)
	r.GET("/cors", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"origins": cfg.GetCORSOrigins(),
			"methods": cfg.GetCORSMethods(),
			"headers": cfg.GetCORSHeaders(),
			"message": "CORS configuration",
		})
	})

	// Public auth routes (no middleware)
	authRoutes := r.Group("/api/auth")
	routes.SetupAuthRoutes(authRoutes)

	// Protected API routes (requires authentication)
	apiRoutes := r.Group("/api")
	apiRoutes.Use(middleware.AuthMiddleware())
	routes.SetupFileRoutes(apiRoutes)
	routes.SetupFolderRoutes(apiRoutes)
	routes.SetupFavoriteRoutes(apiRoutes)
	routes.SetupRecentFilesRoutes(apiRoutes)
	routes.SetupProtectedAuthRoutes(apiRoutes)
	routes.SetupAnalyticsRoutes(apiRoutes)
	routes.SetupSharingRoutes(apiRoutes)

	adminRoutes := r.Group("/api/admin")
	adminRoutes.Use(middleware.AuthMiddleware())
	adminRoutes.Use(middleware.AdminMiddleware())
	routes.SetupAdminRoutes(adminRoutes)

	// Public sharing routes (no authentication required)
	routes.SetupPublicSharingRoutes(r)

	log.Printf("Server starting on port %s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}