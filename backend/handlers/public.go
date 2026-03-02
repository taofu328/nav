package handlers

import (
	"nav-backend/database"
	"nav-backend/models"

	"github.com/gin-gonic/gin"
)

func GetPublicCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取分类失败"})
		return
	}
	c.JSON(200, categories)
}

func GetPublicBookmarks(c *gin.Context) {
	var bookmarks []models.Bookmark
	if err := database.DB.
		Preload("Category").
		Order("sort_order ASC, id ASC").
		Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "获取网址失败"})
		return
	}
	c.JSON(200, bookmarks)
}
