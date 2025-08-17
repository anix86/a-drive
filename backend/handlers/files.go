package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

func ListFiles(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Query("folder_id")
	
	var folders []models.Folder
	var files []models.File
	
	folderQuery := db.Where("user_id = ?", userID)
	fileQuery := db.Where("user_id = ?", userID)
	
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
		"folders": folders,
		"files":   files,
	})
}

func UploadFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	defer file.Close()
	
	folderIDStr := c.PostForm("folder_id")
	var folderID *uint
	if folderIDStr != "" && folderIDStr != "root" {
		fID, err := strconv.ParseUint(folderIDStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid folder ID"})
			return
		}
		uid := uint(fID)
		folderID = &uid
	}
	
	userDir := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID))
	fileName := fmt.Sprintf("%d_%s", userID, header.Filename)
	filePath := filepath.Join(userDir, fileName)
	
	out, err := os.Create(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create file"})
		return
	}
	defer out.Close()
	
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}
	
	fileModel := models.File{
		Name:         header.Filename,
		OriginalName: header.Filename,
		FolderID:     folderID,
		UserID:       userID,
		FilePath:     filePath,
		Size:         header.Size,
		MimeType:     header.Header.Get("Content-Type"),
	}
	
	if err := db.Create(&fileModel).Error; err != nil {
		os.Remove(filePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"file": fileModel})
}

func DownloadFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", file.OriginalName))
	c.File(file.FilePath)
}

func DeleteFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	if err := os.Remove(file.FilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from disk"})
		return
	}
	
	if err := db.Delete(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from database"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func RenameFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	file.Name = req.Name
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename file"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"file": file})
}