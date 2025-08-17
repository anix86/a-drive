package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

type AnalyticsData struct {
	TotalUsers       int64                 `json:"total_users"`
	TotalFiles       int64                 `json:"total_files"`
	TotalFolders     int64                 `json:"total_folders"`
	TotalStorage     int64                 `json:"total_storage"`
	RecentUploads    []models.File         `json:"recent_uploads"`
	TopFileTypes     []FileTypeStats       `json:"top_file_types"`
	StorageByUser    []UserStorageStats    `json:"storage_by_user"`
	ActivityTimeline []ActivityStats       `json:"activity_timeline"`
}

type FileTypeStats struct {
	MimeType string `json:"mime_type"`
	Count    int64  `json:"count"`
	Size     int64  `json:"size"`
}

type UserStorageStats struct {
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
	FileCount int64 `json:"file_count"`
	TotalSize int64 `json:"total_size"`
}

type ActivityStats struct {
	Date      string `json:"date"`
	Uploads   int64  `json:"uploads"`
	Downloads int64  `json:"downloads"`
}

func GetSystemAnalytics(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userRole := c.MustGet("user_role").(string)

	// Only admins can access system analytics
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	analytics := AnalyticsData{}

	// Get total counts
	db.Model(&models.User{}).Count(&analytics.TotalUsers)
	db.Model(&models.File{}).Count(&analytics.TotalFiles)
	db.Model(&models.Folder{}).Count(&analytics.TotalFolders)

	// Get total storage
	var totalStorage struct {
		Total int64
	}
	db.Model(&models.File{}).Select("COALESCE(SUM(size), 0) as total").Scan(&totalStorage)
	analytics.TotalStorage = totalStorage.Total

	// Get recent uploads (last 10)
	db.Preload("User").Order("created_at DESC").Limit(10).Find(&analytics.RecentUploads)

	// Get top file types
	db.Model(&models.File{}).
		Select("mime_type, COUNT(*) as count, COALESCE(SUM(size), 0) as size").
		Group("mime_type").
		Order("count DESC").
		Limit(10).
		Scan(&analytics.TopFileTypes)

	// Get storage by user
	db.Model(&models.File{}).
		Select("files.user_id, users.username, COUNT(*) as file_count, COALESCE(SUM(files.size), 0) as total_size").
		Joins("LEFT JOIN users ON users.id = files.user_id").
		Group("files.user_id, users.username").
		Order("total_size DESC").
		Limit(10).
		Scan(&analytics.StorageByUser)

	// Get activity timeline (last 7 days)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		var uploads int64
		db.Model(&models.File{}).
			Where("DATE(created_at) = ?", dateStr).
			Count(&uploads)

		analytics.ActivityTimeline = append(analytics.ActivityTimeline, ActivityStats{
			Date:    dateStr,
			Uploads: uploads,
			Downloads: 0, // Note: We don't track downloads in the current schema
		})
	}

	c.JSON(http.StatusOK, analytics)
}

type DailyUsage struct {
	Date string `json:"date"`
	Size int64  `json:"size"`
}

func GetUserAnalytics(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	type UserAnalytics struct {
		TotalFiles       int64           `json:"total_files"`
		TotalFolders     int64           `json:"total_folders"`
		TotalStorage     int64           `json:"total_storage"`
		RecentFiles      []models.File   `json:"recent_files"`
		FilesByType      []FileTypeStats `json:"files_by_type"`
		StorageUsage     []DailyUsage    `json:"storage_usage"`
	}

	userAnalytics := UserAnalytics{}

	// Get user's file and folder counts
	db.Model(&models.File{}).Where("user_id = ?", userID).Count(&userAnalytics.TotalFiles)
	db.Model(&models.Folder{}).Where("user_id = ?", userID).Count(&userAnalytics.TotalFolders)

	// Get user's total storage
	var totalStorage struct {
		Total int64
	}
	db.Model(&models.File{}).
		Select("COALESCE(SUM(size), 0) as total").
		Where("user_id = ?", userID).
		Scan(&totalStorage)
	userAnalytics.TotalStorage = totalStorage.Total

	// Get recent files (last 10)
	db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(10).
		Find(&userAnalytics.RecentFiles)

	// Get files by type for this user
	db.Model(&models.File{}).
		Select("mime_type, COUNT(*) as count, COALESCE(SUM(size), 0) as size").
		Where("user_id = ?", userID).
		Group("mime_type").
		Order("count DESC").
		Limit(10).
		Scan(&userAnalytics.FilesByType)

	// Get storage usage over time (last 7 days)
	for i := 6; i >= 0; i-- {
		date := time.Now().AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")

		var size struct {
			Total int64
		}
		db.Model(&models.File{}).
			Select("COALESCE(SUM(size), 0) as total").
			Where("user_id = ? AND DATE(created_at) <= ?", userID, dateStr).
			Scan(&size)

		userAnalytics.StorageUsage = append(userAnalytics.StorageUsage, DailyUsage{
			Date: dateStr,
			Size: size.Total,
		})
	}

	c.JSON(http.StatusOK, userAnalytics)
}