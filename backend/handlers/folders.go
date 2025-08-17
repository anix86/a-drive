package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

func CreateFolder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	var req struct {
		Name     string `json:"name" binding:"required"`
		ParentID *uint  `json:"parent_id"`
		IconType string `json:"icon_type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.IconType == "" {
		req.IconType = "folder"
	}
	
	var path string
	if req.ParentID != nil {
		var parent models.Folder
		if err := db.Where("id = ? AND user_id = ?", *req.ParentID, userID).First(&parent).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Parent folder not found"})
			return
		}
		path = filepath.Join(parent.Path, req.Name)
	} else {
		path = req.Name
	}
	
	folder := models.Folder{
		Name:     req.Name,
		ParentID: req.ParentID,
		UserID:   userID,
		IconType: req.IconType,
		Path:     path,
	}
	
	if err := db.Create(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create folder"})
		return
	}
	
	physicalPath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), path)
	if err := os.MkdirAll(physicalPath, 0755); err != nil {
		db.Delete(&folder)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create physical folder"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{"folder": folder})
}

func GetFolder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Param("id")
	
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).
		Preload("Subfolders").
		Preload("Files").
		First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func GetFolderBreadcrumbs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Param("id")
	
	type BreadcrumbItem struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
	
	var breadcrumbs []BreadcrumbItem
	
	// Always start with Home
	breadcrumbs = append(breadcrumbs, BreadcrumbItem{ID: 0, Name: "Home"})
	
	if folderID == "root" || folderID == "0" {
		c.JSON(http.StatusOK, gin.H{"breadcrumbs": breadcrumbs})
		return
	}
	
	// Get the current folder and build path
	var currentFolder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&currentFolder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}
	
	// Build breadcrumbs by traversing up the parent chain
	var pathFolders []models.Folder
	current := currentFolder
	
	for {
		pathFolders = append([]models.Folder{current}, pathFolders...)
		
		if current.ParentID == nil {
			break
		}
		
		var parent models.Folder
		if err := db.Where("id = ? AND user_id = ?", *current.ParentID, userID).First(&parent).Error; err != nil {
			break
		}
		current = parent
	}
	
	// Convert to breadcrumb items
	for _, folder := range pathFolders {
		breadcrumbs = append(breadcrumbs, BreadcrumbItem{
			ID:   folder.ID,
			Name: folder.Name,
		})
	}
	
	c.JSON(http.StatusOK, gin.H{"breadcrumbs": breadcrumbs})
}

func UpdateFolder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Param("id")
	
	var req struct {
		Name     string `json:"name"`
		IconType string `json:"icon_type"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}
	
	if req.Name != "" && req.Name != folder.Name {
		oldPath := folder.Path
		newPath := filepath.Join(filepath.Dir(folder.Path), req.Name)
		
		oldPhysicalPath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), oldPath)
		newPhysicalPath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), newPath)
		
		if err := os.Rename(oldPhysicalPath, newPhysicalPath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to rename physical folder"})
			return
		}
		
		folder.Name = req.Name
		folder.Path = newPath
	}
	
	if req.IconType != "" {
		folder.IconType = req.IconType
	}
	
	if err := db.Save(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update folder"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"folder": folder})
}

func DeleteFolder(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Param("id")
	
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}
	
	physicalPath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), folder.Path)
	if err := os.RemoveAll(physicalPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete physical folder"})
		return
	}
	
	if err := db.Delete(&folder).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete folder from database"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Folder deleted successfully"})
}

func CreateZipArchive(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	folderID := c.Param("id")
	
	var folder models.Folder
	if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
		return
	}
	
	tempDir := os.TempDir()
	zipFileName := fmt.Sprintf("%s_%d.zip", strings.ReplaceAll(folder.Name, " ", "_"), userID)
	zipPath := filepath.Join(tempDir, zipFileName)
	
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file"})
		return
	}
	defer zipFile.Close()
	defer os.Remove(zipPath)
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	sourcePath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), folder.Path)
	
	err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, err := filepath.Rel(sourcePath, path)
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			return nil
		}
		
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		
		writer, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}
		
		_, err = io.Copy(writer, file)
		return err
	})
	
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip archive"})
		return
	}
	
	zipWriter.Close()
	zipFile.Close()
	
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	c.Header("Content-Type", "application/zip")
	c.File(zipPath)
}