package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
	"a-drive-backend/utils"
)

func ListUsers(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var users []models.User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func CreateUser(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var existingUser models.User
	if err := db.Where("username = ? OR email = ?", req.Username, req.Email).First(&existingUser).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}
	
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	
	user := models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}
	
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	
	userDir := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", user.ID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user directory"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

func BrowseUserFiles(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	userIDParam := c.Query("user_id")
	if userIDParam == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id parameter required"})
		return
	}
	
	userID, err := strconv.ParseUint(userIDParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	
	var user models.User
	if err := db.First(&user, uint(userID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	
	folderID := c.Query("folder_id")
	
	var folders []models.Folder
	var files []models.File
	
	folderQuery := db.Where("user_id = ?", uint(userID))
	fileQuery := db.Where("user_id = ?", uint(userID))
	
	if folderID != "" {
		if folderID == "root" {
			folderQuery = folderQuery.Where("parent_id IS NULL")
			fileQuery = fileQuery.Where("folder_id IS NULL")
		} else {
			fID, err := strconv.ParseUint(folderID, 10, 32)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
				return
			}
			folderQuery = folderQuery.Where("parent_id = ?", uint(fID))
			fileQuery = fileQuery.Where("folder_id = ?", uint(fID))
		}
	}
	
	if err := folderQuery.Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch folders"})
		return
	}
	
	if err := fileQuery.Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch files"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"user":    user,
		"folders": folders,
		"files":   files,
	})
}