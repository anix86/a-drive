package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"a-drive-backend/models"
	"a-drive-backend/utils"
)

func Init(databasePath string) *gorm.DB {
	dir := filepath.Dir(databasePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		log.Fatal("Failed to create database directory:", err)
	}

	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := migrate(db); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	if err := createAdminUser(db); err != nil {
		log.Println("Admin user creation skipped:", err)
	}

	return db
}

func migrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Folder{},
		&models.File{},
		&models.FileVersion{},
		&models.FileShare{},
		&models.ShareAccess{},
		&models.Favorite{},
		&models.RecentAccess{},
	)
}

func createAdminUser(db *gorm.DB) error {
	var count int64
	db.Model(&models.User{}).Count(&count)
	
	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	adminUser := models.User{
		Username:     "admin",
		Email:        "admin@example.com",
		PasswordHash: hashedPassword,
		Role:         "admin",
	}

	if err := db.Create(&adminUser).Error; err != nil {
		return err
	}

	userDir := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", adminUser.ID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		log.Printf("Failed to create admin user directory: %v", err)
	}

	log.Println("Created admin user - Username: admin, Password: admin123")
	return nil
}