package main

import (
	"log"
	"os"

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
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.DatabaseMiddleware(db))

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "OK", "message": "A-Drive is running"})
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}