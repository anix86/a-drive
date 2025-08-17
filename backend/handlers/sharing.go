package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

type CreateShareRequest struct {
	ShareType     string     `json:"share_type" binding:"required,oneof=public password private"`
	Password      string     `json:"password"`
	ExpiresAt     *time.Time `json:"expires_at"`
	MaxDownloads  *int       `json:"max_downloads"`
	AllowPreview  bool       `json:"allow_preview"`
}

type ShareAccessRequest struct {
	Password string `json:"password"`
}

func CreateFileShare(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	var req CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Verify file exists and belongs to user
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	// Generate unique share token
	token, err := generateShareToken()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate share token"})
		return
	}
	
	// Create share record
	share := models.FileShare{
		FileID:       file.ID,
		SharedBy:     userID,
		ShareToken:   token,
		ShareType:    req.ShareType,
		ExpiresAt:    req.ExpiresAt,
		MaxDownloads: req.MaxDownloads,
		AllowPreview: req.AllowPreview,
	}
	
	// Hash password if provided
	if req.ShareType == "password" && req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		share.Password = string(hashedPassword)
	}
	
	if err := db.Create(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create share"})
		return
	}
	
	// Load the created share with relationships
	db.Preload("File").Preload("SharedByUser").First(&share, share.ID)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "File shared successfully",
		"share":   share,
		"share_url": getShareURL(c, token),
	})
}

func GetFileShares(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	fileID := c.Param("id")
	
	// Verify file exists and belongs to user
	var file models.File
	if err := db.Where("id = ? AND user_id = ?", fileID, userID).First(&file).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}
	
	var shares []models.FileShare
	if err := db.Where("file_id = ?", file.ID).
		Preload("SharedByUser").
		Order("created_at DESC").
		Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shares"})
		return
	}
	
	// Add share URLs
	for i := range shares {
		shares[i].ShareToken = getShareURL(c, shares[i].ShareToken)
	}
	
	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func GetUserShares(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	var shares []models.FileShare
	if err := db.Where("shared_by = ?", userID).
		Preload("File").
		Preload("SharedByUser").
		Order("created_at DESC").
		Find(&shares).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch shares"})
		return
	}
	
	// Add share URLs
	for i := range shares {
		shares[i].ShareToken = getShareURL(c, shares[i].ShareToken)
	}
	
	c.JSON(http.StatusOK, gin.H{"shares": shares})
}

func DeleteFileShare(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)
	
	shareID := c.Param("share_id")
	
	var share models.FileShare
	if err := db.Where("id = ? AND shared_by = ?", shareID, userID).First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}
	
	if err := db.Delete(&share).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete share"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Share deleted successfully"})
}

// Public share access endpoints (no auth required)
func AccessSharedFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	token := c.Param("token")
	
	var share models.FileShare
	if err := db.Where("share_token = ?", token).
		Preload("File").
		Preload("SharedByUser").
		First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}
	
	// Check if share has expired
	if share.ExpiresAt != nil && time.Now().After(*share.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Share has expired"})
		return
	}
	
	// Check download limit
	if share.MaxDownloads != nil && share.DownloadCount >= *share.MaxDownloads {
		c.JSON(http.StatusGone, gin.H{"error": "Download limit exceeded"})
		return
	}
	
	// For password-protected shares, require password verification
	if share.ShareType == "password" {
		var req ShareAccessRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Password required"})
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(share.Password), []byte(req.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}
	}
	
	// Log the access
	logShareAccess(db, share.ID, c.ClientIP(), c.GetHeader("User-Agent"), "view")
	
	// Return share info (without sensitive data)
	shareInfo := gin.H{
		"id":            share.ID,
		"file":          share.File,
		"shared_by":     share.SharedByUser.Username,
		"share_type":    share.ShareType,
		"allow_preview": share.AllowPreview,
		"created_at":    share.CreatedAt,
	}
	
	c.JSON(http.StatusOK, shareInfo)
}

func DownloadSharedFile(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	
	token := c.Param("token")
	
	var share models.FileShare
	if err := db.Where("share_token = ?", token).
		Preload("File").
		First(&share).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Share not found"})
		return
	}
	
	// Check if share has expired
	if share.ExpiresAt != nil && time.Now().After(*share.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "Share has expired"})
		return
	}
	
	// Check download limit
	if share.MaxDownloads != nil && share.DownloadCount >= *share.MaxDownloads {
		c.JSON(http.StatusGone, gin.H{"error": "Download limit exceeded"})
		return
	}
	
	// For password-protected shares, password should be verified in previous step
	// This endpoint assumes the user has already been authenticated via AccessSharedFile
	
	// Increment download count
	db.Model(&share).Update("download_count", gorm.Expr("download_count + 1"))
	
	// Log the download
	logShareAccess(db, share.ID, c.ClientIP(), c.GetHeader("User-Agent"), "download")
	
	// Serve the file
	c.FileAttachment(share.File.FilePath, share.File.Name)
}

// Helper functions
func generateShareToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func getShareURL(c *gin.Context, token string) string {
	scheme := "http"
	if c.GetHeader("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	
	host := c.GetHeader("Host")
	if host == "" {
		host = "localhost:3000" // fallback for development
	}
	
	return scheme + "://" + host + "/share/" + token
}

func logShareAccess(db *gorm.DB, shareID uint, ip, userAgent, action string) {
	access := models.ShareAccess{
		ShareID:    shareID,
		AccessedBy: ip,
		AccessedAt: time.Now(),
		UserAgent:  userAgent,
		Action:     action,
	}
	db.Create(&access)
}