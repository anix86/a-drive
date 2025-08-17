package handlers

import (
	"crypto/md5"
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

func EnableVersioning(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	if file.VersioningEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Versioning already enabled for this file"})
		return
	}
	
	// Create initial version entry
	checksum, err := calculateFileChecksum(file.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate file checksum"})
		return
	}
	
	version := models.FileVersion{
		FileID:    file.ID,
		Version:   1,
		FilePath:  file.FilePath,
		Size:      file.Size,
		Checksum:  checksum,
		Comment:   "Initial version",
		CreatedBy: userID,
	}
	
	if err := db.Create(&version).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create version record"})
		return
	}
	
	// Enable versioning on the file
	file.VersioningEnabled = true
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to enable versioning"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Versioning enabled successfully"})
}

func DisableVersioning(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	// Get all versions except the current one
	var oldVersions []models.FileVersion
	if err := db.Where("file_id = ? AND version != ?", file.ID, file.CurrentVersion).Find(&oldVersions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions"})
		return
	}
	
	// Delete old version files
	for _, version := range oldVersions {
		if version.FilePath != file.FilePath {
			os.Remove(version.FilePath)
		}
	}
	
	// Delete version records
	if err := db.Where("file_id = ?", file.ID).Delete(&models.FileVersion{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete version records"})
		return
	}
	
	// Disable versioning
	file.VersioningEnabled = false
	file.CurrentVersion = 1
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to disable versioning"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Versioning disabled successfully"})
}

func GetFileVersions(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	if !file.VersioningEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Versioning not enabled for this file"})
		return
	}
	
	var versions []models.FileVersion
	if err := db.Where("file_id = ?", file.ID).
		Preload("CreatedByUser").
		Order("version DESC").
		Find(&versions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch versions"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"file":     file,
		"versions": versions,
	})
}

func CreateNewVersion(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	comment := c.PostForm("comment")
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	if !file.VersioningEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Versioning not enabled for this file"})
		return
	}
	
	uploadedFile, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}
	
	// Open uploaded file
	src, err := uploadedFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer src.Close()
	
	// Create new version file path
	newVersion := file.CurrentVersion + 1
	fileExt := filepath.Ext(file.FilePath)
	baseName := file.FilePath[:len(file.FilePath)-len(fileExt)]
	
	// Copy current file to versioned location
	oldVersionPath := fmt.Sprintf("%s_v%d%s", baseName, file.CurrentVersion, fileExt)
	if err := copyFile(file.FilePath, oldVersionPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to backup current version"})
		return
	}
	
	// Create version record for the old file
	oldChecksum, _ := calculateFileChecksum(file.FilePath)
	oldVersion := models.FileVersion{
		FileID:    file.ID,
		Version:   file.CurrentVersion,
		FilePath:  oldVersionPath,
		Size:      file.Size,
		Checksum:  oldChecksum,
		Comment:   "Previous version",
		CreatedBy: userID,
	}
	
	if err := db.Create(&oldVersion).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create old version record"})
		return
	}
	
	// Save new file
	dst, err := os.Create(file.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new file"})
		return
	}
	defer dst.Close()
	
	size, err := io.Copy(dst, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save new file"})
		return
	}
	
	// Calculate new checksum
	newChecksum, err := calculateFileChecksum(file.FilePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate new file checksum"})
		return
	}
	
	// Create new version record
	newVersionRecord := models.FileVersion{
		FileID:    file.ID,
		Version:   newVersion,
		FilePath:  file.FilePath,
		Size:      size,
		Checksum:  newChecksum,
		Comment:   comment,
		CreatedBy: userID,
	}
	
	if err := db.Create(&newVersionRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create new version record"})
		return
	}
	
	// Update file record
	file.CurrentVersion = newVersion
	file.Size = size
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file record"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "New version created successfully",
		"version": newVersionRecord,
	})
}

func RestoreVersion(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	versionIDStr := c.Param("version_id")
	
	versionID, err := strconv.ParseUint(versionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	if !file.VersioningEnabled {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Versioning not enabled for this file"})
		return
	}
	
	var version models.FileVersion
	if err := db.Where("id = ? AND file_id = ?", versionID, file.ID).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}
	
	// Copy version file to current location
	if err := copyFile(version.FilePath, file.FilePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to restore version"})
		return
	}
	
	// Update file record
	file.Size = version.Size
	if err := db.Save(&file).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update file record"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Version restored successfully",
		"restored_version": version.Version,
	})
}

func DownloadVersion(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	versionIDStr := c.Param("version_id")
	
	versionID, err := strconv.ParseUint(versionIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid version ID"})
		return
	}
	
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	var version models.FileVersion
	if err := db.Where("id = ? AND file_id = ?", versionID, file.ID).First(&version).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Version not found"})
		return
	}
	
	c.FileAttachment(version.FilePath, fmt.Sprintf("%s_v%d%s", 
		file.Name[:len(file.Name)-len(filepath.Ext(file.Name))], 
		version.Version, 
		filepath.Ext(file.Name)))
}

// Helper functions
func calculateFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()
	
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	
	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()
	
	_, err = io.Copy(destFile, sourceFile)
	return err
}