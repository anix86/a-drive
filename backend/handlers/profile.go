package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
	"a-drive-backend/utils"
)

type UpdateProfileRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required,min=6"`
}

type UserStats struct {
	TotalFiles   int   `json:"total_files"`
	TotalFolders int   `json:"total_folders"`
	TotalSize    int64 `json:"total_size"`
	LastActivity string `json:"last_activity"`
}

func GetProfile(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	db := c.MustGet("db").(*gorm.DB)
	
	var stats UserStats
	
	// Get file count and total size
	var fileCount, folderCount int64
	db.Model(&models.File{}).Where("user_id = ?", user.ID).Count(&fileCount)
	db.Model(&models.Folder{}).Where("user_id = ?", user.ID).Count(&folderCount)
	stats.TotalFiles = int(fileCount)
	stats.TotalFolders = int(folderCount)
	db.Model(&models.File{}).Where("user_id = ?", user.ID).Select("COALESCE(SUM(size), 0)").Scan(&stats.TotalSize)
	
	// Get last activity (most recent file upload or folder creation)
	var lastFile models.File
	var lastFolder models.Folder
	
	db.Where("user_id = ?", user.ID).Order("created_at DESC").First(&lastFile)
	db.Where("user_id = ?", user.ID).Order("created_at DESC").First(&lastFolder)
	
	if lastFile.CreatedAt.After(lastFolder.CreatedAt) {
		stats.LastActivity = lastFile.CreatedAt.Format("2006-01-02 15:04:05")
	} else if !lastFolder.CreatedAt.IsZero() {
		stats.LastActivity = lastFolder.CreatedAt.Format("2006-01-02 15:04:05")
	} else {
		stats.LastActivity = "No activity yet"
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
			"created_at": user.CreatedAt,
		},
		"stats": stats,
	})
}

func UpdateProfile(c *gin.Context) {
	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user := c.MustGet("user").(models.User)
	db := c.MustGet("db").(*gorm.DB)
	
	// Check if username or email already exists (if changed)
	if req.Username != "" && req.Username != user.Username {
		var existingUser models.User
		if err := db.Where("username = ? AND id != ?", req.Username, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Username already exists"})
			return
		}
		user.Username = req.Username
	}
	
	if req.Email != "" && req.Email != user.Email {
		var existingUser models.User
		if err := db.Where("email = ? AND id != ?", req.Email, user.ID).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
			return
		}
		user.Email = req.Email
	}
	
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"role":     user.Role,
		},
	})
}

func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	user := c.MustGet("user").(models.User)
	db := c.MustGet("db").(*gorm.DB)
	
	// Verify current password
	if !utils.CheckPasswordHash(req.CurrentPassword, user.PasswordHash) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Current password is incorrect"})
		return
	}
	
	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}
	
	// Update password
	user.PasswordHash = hashedPassword
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}