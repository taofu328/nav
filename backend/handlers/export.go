package handlers

import (
	"nav-backend/database"
	"nav-backend/models"

	"github.com/gin-gonic/gin"
)

type ExportDataStruct struct {
	Categories []models.Category `json:"categories"`
	Bookmarks  []models.Bookmark `json:"bookmarks"`
}

func ExportData(c *gin.Context) {
	userID := c.GetUint("user_id")

	var categories []models.Category
	if err := database.DB.Where("user_id = ?", userID).Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to export categories"})
		return
	}

	var bookmarks []models.Bookmark
	if err := database.DB.Where("user_id = ?", userID).Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to export bookmarks"})
		return
	}

	data := ExportDataStruct{
		Categories: categories,
		Bookmarks:  bookmarks,
	}

	c.JSON(200, data)
}

func ImportData(c *gin.Context) {
	userID := c.GetUint("user_id")

	var data ExportDataStruct
	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()

	for _, category := range data.Categories {
		category.ID = 0
		category.UserID = userID
		if err := tx.Create(&category).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to import categories"})
			return
		}
	}

	categoryMap := make(map[uint]*uint)
	var importedCategories []models.Category
	if err := tx.Where("user_id = ?", userID).Find(&importedCategories).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to get imported categories"})
		return
	}

	for i, category := range importedCategories {
		if i < len(data.Categories) {
			categoryID := category.ID
			categoryMap[data.Categories[i].ID] = &categoryID
		}
	}

	for _, bookmark := range data.Bookmarks {
		bookmark.ID = 0
		bookmark.UserID = userID
		if bookmark.CategoryID != nil {
			if newCategoryID, exists := categoryMap[*bookmark.CategoryID]; exists {
				bookmark.CategoryID = newCategoryID
			}
		}
		if err := tx.Create(&bookmark).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "Failed to import bookmarks"})
			return
		}
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Data imported successfully"})
}
