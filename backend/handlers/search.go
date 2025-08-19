package handlers

import (
	"net/http"
	"regexp"
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

// processWildcardQuery converts wildcard patterns to SQL LIKE patterns
func processWildcardQuery(query string) (string, bool) {
	// Check if query contains wildcard characters
	if !strings.Contains(query, "*") && !strings.Contains(query, "?") {
		return query, false
	}
	
	// Escape existing SQL special characters
	escaped := strings.ReplaceAll(query, "%", "\\%")
	escaped = strings.ReplaceAll(escaped, "_", "\\_")
	
	// Convert wildcard patterns to SQL LIKE patterns
	// * matches any sequence of characters (including empty)
	// ? matches any single character
	escaped = strings.ReplaceAll(escaped, "*", "%")
	escaped = strings.ReplaceAll(escaped, "?", "_")
	
	return escaped, true
}

// isExtensionPattern checks if query is an extension pattern like *.pdf
func isExtensionPattern(query string) (string, bool) {
	// Match patterns like *.ext
	re := regexp.MustCompile(`^\*\.([a-zA-Z0-9]+)$`)
	matches := re.FindStringSubmatch(query)
	if len(matches) == 2 {
		return strings.ToLower(matches[1]), true
	}
	return "", false
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
	
	// Process wildcard patterns
	processedQuery, isWildcard := processWildcardQuery(query)
	
	// Check if it's an extension pattern like *.pdf
	extension, isExtensionPattern := isExtensionPattern(query)
	
	var searchPattern string
	if isWildcard {
		searchPattern = processedQuery
	} else {
		searchPattern = "%" + query + "%"
	}
	
	
	// Search files
	fileQuery := db.Where("user_id = ?", userID)
	
	if isExtensionPattern {
		// For extension patterns like *.pdf, search by name ending with extension
		fileQuery = fileQuery.Where("name LIKE ?", "%."+extension)
	} else {
		// Regular name search with wildcard support
		fileQuery = fileQuery.Where("name LIKE ?", searchPattern)
	}
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
	
	// Search folders (only if not an extension pattern, since folders don't have extensions)
	if !isExtensionPattern {
		folderQuery := db.Where("user_id = ? AND name LIKE ?", userID, searchPattern)
		if err := folderQuery.Limit(25).Find(&folders).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search folders"})
			return
		}
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