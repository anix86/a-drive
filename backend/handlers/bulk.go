package handlers

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

type BulkOperationRequest struct {
	FileIDs   []uint `json:"file_ids"`
	FolderIDs []uint `json:"folder_ids"`
	Action    string `json:"action"` // "delete", "move", "download"
	TargetID  *uint  `json:"target_id,omitempty"` // for move operation
}

type BulkOperationResult struct {
	Success     bool     `json:"success"`
	Message     string   `json:"message"`
	Processed   int      `json:"processed"`
	Failed      int      `json:"failed"`
	FailedItems []string `json:"failed_items,omitempty"`
}

func BulkOperation(c *gin.Context) {
	var req BulkOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	if req.Action == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Action is required"})
		return
	}
	
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	result := BulkOperationResult{}
	
	switch req.Action {
	case "delete":
		result = bulkDelete(db, userID, req.FileIDs, req.FolderIDs)
	case "move":
		if req.TargetID == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Target folder ID required for move operation"})
			return
		}
		result = bulkMove(db, userID, req.FileIDs, req.FolderIDs, *req.TargetID)
	case "download":
		// For download, create a ZIP file with all selected items
		result = bulkDownload(c, db, userID, req.FileIDs, req.FolderIDs)
		return // bulkDownload handles the response
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action"})
		return
	}
	
	c.JSON(http.StatusOK, result)
}

func bulkDelete(db *gorm.DB, userID uint, fileIDs, folderIDs []uint) BulkOperationResult {
	result := BulkOperationResult{Success: true}
	
	// Delete files
	for _, fileID := range fileIDs {
		var file models.File
		if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, fmt.Sprintf("file_%d", fileID))
			continue
		}
		
		// Delete physical file
		if err := os.Remove(file.FilePath); err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, file.Name)
			continue
		}
		
		// Delete from database
		if err := db.Delete(&file).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, file.Name)
			continue
		}
		
		result.Processed++
	}
	
	// Delete folders
	for _, folderID := range folderIDs {
		var folder models.Folder
		if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, fmt.Sprintf("folder_%d", folderID))
			continue
		}
		
		// Delete physical folder
		physicalPath := filepath.Join(os.Getenv("ROOT_DIRECTORY"), "root", fmt.Sprintf("%d", userID), folder.Path)
		if err := os.RemoveAll(physicalPath); err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, folder.Name)
			continue
		}
		
		// Delete from database (cascade will handle subfolders and files)
		if err := db.Delete(&folder).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, folder.Name)
			continue
		}
		
		result.Processed++
	}
	
	result.Message = fmt.Sprintf("Processed %d items, %d failed", result.Processed, result.Failed)
	if result.Failed > 0 {
		result.Success = false
	}
	
	return result
}

func bulkMove(db *gorm.DB, userID uint, fileIDs, folderIDs []uint, targetFolderID uint) BulkOperationResult {
	result := BulkOperationResult{Success: true}
	
	// Verify target folder exists and belongs to user
	var targetFolder models.Folder
	if targetFolderID > 0 {
		if err := db.Where("id = ? AND user_id = ?", targetFolderID, userID).First(&targetFolder).Error; err != nil {
			result.Success = false
			result.Message = "Target folder not found"
			return result
		}
	}
	
	// Move files
	for _, fileID := range fileIDs {
		var file models.File
		if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, fmt.Sprintf("file_%d", fileID))
			continue
		}
		
		// Update folder_id
		if targetFolderID == 0 {
			file.FolderID = nil // Move to root
		} else {
			file.FolderID = &targetFolderID
		}
		
		if err := db.Save(&file).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, file.Name)
			continue
		}
		
		result.Processed++
	}
	
	// Move folders
	for _, folderID := range folderIDs {
		var folder models.Folder
		if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, fmt.Sprintf("folder_%d", folderID))
			continue
		}
		
		// Update parent_id
		if targetFolderID == 0 {
			folder.ParentID = nil // Move to root
			folder.Path = folder.Name
		} else {
			folder.ParentID = &targetFolderID
			folder.Path = filepath.Join(targetFolder.Path, folder.Name)
		}
		
		if err := db.Save(&folder).Error; err != nil {
			result.Failed++
			result.FailedItems = append(result.FailedItems, folder.Name)
			continue
		}
		
		result.Processed++
	}
	
	result.Message = fmt.Sprintf("Moved %d items, %d failed", result.Processed, result.Failed)
	if result.Failed > 0 {
		result.Success = false
	}
	
	return result
}

func bulkDownload(c *gin.Context, db *gorm.DB, userID uint, fileIDs, folderIDs []uint) BulkOperationResult {
	result := BulkOperationResult{Success: true}
	
	// Create temporary ZIP file
	tempDir := os.TempDir()
	zipFileName := fmt.Sprintf("bulk_download_%d.zip", userID)
	zipPath := filepath.Join(tempDir, zipFileName)
	
	zipFile, err := os.Create(zipPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create zip file"})
		return result
	}
	defer zipFile.Close()
	defer os.Remove(zipPath)
	
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()
	
	// Add files to ZIP
	for _, fileID := range fileIDs {
		var file models.File
		if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
			continue
		}
		
		fileReader, err := os.Open(file.FilePath)
		if err != nil {
			continue
		}
		defer fileReader.Close()
		
		writer, err := zipWriter.Create(file.Name)
		if err != nil {
			continue
		}
		
		_, err = io.Copy(writer, fileReader)
		if err != nil {
			continue
		}
		
		result.Processed++
	}
	
	// Add folders to ZIP
	for _, folderID := range folderIDs {
		var folder models.Folder
		if err := db.Where("id = ? AND user_id = ?", folderID, userID).First(&folder).Error; err != nil {
			continue
		}
		
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
			
			zipPath := filepath.Join(folder.Name, relPath)
			
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			
			writer, err := zipWriter.Create(zipPath)
			if err != nil {
				return err
			}
			
			_, err = io.Copy(writer, file)
			return err
		})
		
		if err == nil {
			result.Processed++
		}
	}
	
	zipWriter.Close()
	zipFile.Close()
	
	// Send ZIP file as response
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	c.Header("Content-Type", "application/zip")
	c.File(zipPath)
	
	result.Success = true
	result.Message = fmt.Sprintf("Created ZIP with %d items", result.Processed)
	return result
}