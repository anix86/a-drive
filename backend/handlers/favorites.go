package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"a-drive-backend/models"
)

// GetFavorites returns all favorites for the authenticated user
func GetFavorites(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	var favorites []models.Favorite
	if err := db.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch favorites"})
		return
	}

	// Populate item details for each favorite
	var populatedFavorites []map[string]interface{}
	for _, fav := range favorites {
		favoriteData := map[string]interface{}{
			"id":         fav.ID,
			"user_id":    fav.UserID,
			"item_type":  fav.ItemType,
			"item_id":    fav.ItemID,
			"created_at": fav.CreatedAt,
		}

		if fav.ItemType == "file" {
			var file models.File
			if err := db.Where("id = ? AND user_id = ?", fav.ItemID, userID).First(&file).Error; err == nil {
				file.IsFavorite = true
				favoriteData["item"] = file
			}
		} else if fav.ItemType == "folder" {
			var folder models.Folder
			if err := db.Where("id = ? AND user_id = ?", fav.ItemID, userID).First(&folder).Error; err == nil {
				folder.IsFavorite = true
				favoriteData["item"] = folder
			}
		}

		populatedFavorites = append(populatedFavorites, favoriteData)
	}

	c.JSON(http.StatusOK, gin.H{"favorites": populatedFavorites})
}

// AddFavorite adds an item to favorites
func AddFavorite(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	var req struct {
		ItemType string `json:"item_type" binding:"required,oneof=file folder"`
		ItemID   uint   `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the item exists and belongs to the user
	if req.ItemType == "file" {
		var file models.File
		if err := db.Where("id = ? AND user_id = ?", req.ItemID, userID).First(&file).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
			return
		}
	} else {
		var folder models.Folder
		if err := db.Where("id = ? AND user_id = ?", req.ItemID, userID).First(&folder).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Folder not found"})
			return
		}
	}

	// Check if already favorited
	var existing models.Favorite
	if err := db.Where("user_id = ? AND item_type = ? AND item_id = ?", userID, req.ItemType, req.ItemID).First(&existing).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Item is already favorited"})
		return
	}

	// Create favorite
	favorite := models.Favorite{
		UserID:   userID,
		ItemType: req.ItemType,
		ItemID:   req.ItemID,
	}

	if err := db.Create(&favorite).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add favorite"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"favorite": favorite})
}

// RemoveFavorite removes an item from favorites
func RemoveFavorite(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	favoriteIDStr := c.Param("id")
	favoriteID, err := strconv.ParseUint(favoriteIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid favorite ID"})
		return
	}

	// Delete the favorite (ensure it belongs to the user)
	result := db.Where("id = ? AND user_id = ?", uint(favoriteID), userID).Delete(&models.Favorite{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}

// CheckFavorite checks if an item is favorited
func CheckFavorite(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	itemType := c.Param("type")
	itemIDStr := c.Param("id")

	if itemType != "file" && itemType != "folder" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item type"})
		return
	}

	itemID, err := strconv.ParseUint(itemIDStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid item ID"})
		return
	}

	var favorite models.Favorite
	if err := db.Where("user_id = ? AND item_type = ? AND item_id = ?", userID, itemType, uint(itemID)).First(&favorite).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{"is_favorite": false})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check favorite status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_favorite": true,
		"favorite_id": favorite.ID,
	})
}

// RemoveFavoriteByItem removes a favorite by item type and ID
func RemoveFavoriteByItem(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	userID := c.MustGet("user_id").(uint)

	var req struct {
		ItemType string `json:"item_type" binding:"required,oneof=file folder"`
		ItemID   uint   `json:"item_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Delete the favorite
	result := db.Where("user_id = ? AND item_type = ? AND item_id = ?", userID, req.ItemType, req.ItemID).Delete(&models.Favorite{})
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove favorite"})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Favorite not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Favorite removed successfully"})
}