package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

// GetRecentFiles returns the most recent 20 files and folders accessed by the user
func GetRecentFiles(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	var recentAccesses []models.RecentAccess
	
	// Get the most recent 20 accesses, ordered by accessed_at descending
	if err := db.Where("user_id = ?", userID).
		Order("accessed_at DESC").
		Limit(20).
		Find(&recentAccesses).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recent files"})
		return
	}

	// For each recent access, fetch the actual file or folder data
	var result []gin.H
	for _, access := range recentAccesses {
		var item gin.H
		
		if access.ItemType == "file" {
			var file models.File
			if err := db.First(&file, access.ItemID).Error; err != nil {
				// File might have been deleted, skip it
				continue
			}
			item = gin.H{
				"id":            access.ID,
				"user_id":       access.UserID,
				"item_type":     access.ItemType,
				"item_id":       access.ItemID,
				"accessed_at":   access.AccessedAt,
				"created_at":    access.CreatedAt,
				"item":          file,
			}
		} else if access.ItemType == "folder" {
			var folder models.Folder
			if err := db.First(&folder, access.ItemID).Error; err != nil {
				// Folder might have been deleted, skip it
				continue
			}
			item = gin.H{
				"id":            access.ID,
				"user_id":       access.UserID,
				"item_type":     access.ItemType,
				"item_id":       access.ItemID,
				"accessed_at":   access.AccessedAt,
				"created_at":    access.CreatedAt,
				"item":          folder,
			}
		}
		
		if item != nil {
			result = append(result, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{"recent_files": result})
}

// TrackFileAccess records when a user accesses a file
func TrackFileAccess(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileIDStr := c.Param("id")
	fileID, err := strconv.ParseUint(fileIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Check if file exists and belongs to user
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if recent access already exists
	var existingAccess models.RecentAccess
	result := db.Where("user_id = ? AND item_type = ? AND item_id = ?", 
		userID, "file", fileID).First(&existingAccess)
	
	if result.Error == nil {
		// Update existing access time
		existingAccess.AccessedAt = time.Now()
		if err := db.Save(&existingAccess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update access time"})
			return
		}
	} else {
		// Create new access record
		recentAccess := models.RecentAccess{
			UserID:     uint(userID),
			ItemType:   "file",
			ItemID:     uint(fileID),
			AccessedAt: time.Now(),
		}
		if err := db.Create(&recentAccess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track access"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access tracked successfully"})
}

// TrackFolderAccess records when a user accesses a folder
func TrackFolderAccess(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderIDStr := c.Param("id")
	folderID, err := strconv.ParseUint(folderIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
		return
	}

	// Check if folder exists and belongs to user
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}

	// Check if recent access already exists
	var existingAccess models.RecentAccess
	result := db.Where("user_id = ? AND item_type = ? AND item_id = ?", 
		userID, "folder", folderID).First(&existingAccess)
	
	if result.Error == nil {
		// Update existing access time
		existingAccess.AccessedAt = time.Now()
		if err := db.Save(&existingAccess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update access time"})
			return
		}
	} else {
		// Create new access record
		recentAccess := models.RecentAccess{
			UserID:     uint(userID),
			ItemType:   "folder",
			ItemID:     uint(folderID),
			AccessedAt: time.Now(),
		}
		if err := db.Create(&recentAccess).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to track access"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Access tracked successfully"})
}