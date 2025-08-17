package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

type SearchResult struct {
	Files   []models.File   `json:"files"`
	Folders []models.Folder `json:"folders"`
	Total   int             `json:"total"`
}

func SearchFiles(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	query := c.Query("q")
	fileType := c.Query("type") // optional filter by file type
	
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
		return
	}
	
	var files []models.File
	var folders []models.Folder
	
	// Search files
	fileQuery := db.Where("user_id = ? AND name LIKE ?", userID, "%"+query+"%")
	if fileType != "" {
		switch strings.ToLower(fileType) {
		case "image":
			fileQuery = fileQuery.Where("mime_type LIKE ?", "image/%")
		case "document":
			fileQuery = fileQuery.Where("mime_type IN ?", []string{
				"application/pdf", "application/msword", "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				"text/plain", "text/csv", "application/vnd.ms-excel", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			})
		case "video":
			fileQuery = fileQuery.Where("mime_type LIKE ?", "video/%")
		case "audio":
			fileQuery = fileQuery.Where("mime_type LIKE ?", "audio/%")
		case "archive":
			fileQuery = fileQuery.Where("mime_type IN ?", []string{
				"application/zip", "application/x-rar-compressed", "application/x-7z-compressed",
				"application/gzip", "application/x-tar",
			})
		}
	}
	
	if err := fileQuery.Limit(25).Find(&files).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search files"})
		return
	}
	
	// Search folders
	if err := db.Where("user_id = ? AND name LIKE ?", userID, "%"+query+"%").
		Limit(25).Find(&folders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search folders"})
		return
	}
	
	result := SearchResult{
		Files:   files,
		Folders: folders,
		Total:   len(files) + len(folders),
	}
	
	c.JSON(http.StatusOK, result)
}

func GetFileTypes(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	var mimeTypes []string
	if err := db.Model(&models.File{}).Where("user_id = ?", userID).
		Distinct("mime_type").Pluck("mime_type", &mimeTypes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get file types"})
		return
	}
	
	// Group mime types into categories
	categories := make(map[string][]string)
	for _, mimeType := range mimeTypes {
		if mimeType == "" {
			continue
		}
		
		parts := strings.Split(mimeType, "/")
		if len(parts) >= 1 {
			category := parts[0]
			if category == "application" {
				// Special handling for common application types
				if strings.Contains(mimeType, "pdf") {
					category = "document"
				} else if strings.Contains(mimeType, "zip") || strings.Contains(mimeType, "rar") || 
					strings.Contains(mimeType, "tar") || strings.Contains(mimeType, "gzip") {
					category = "archive"
				} else if strings.Contains(mimeType, "word") || strings.Contains(mimeType, "excel") || 
					strings.Contains(mimeType, "powerpoint") || strings.Contains(mimeType, "spreadsheet") {
					category = "document"
				}
			}
			categories[category] = append(categories[category], mimeType)
		}
	}
	
	c.JSON(http.StatusOK, gin.H{"categories": categories})
}