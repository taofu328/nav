package handlers

import (
	"nav-backend/database"
	"nav-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookmarks(c *gin.Context) {
	var bookmarks []models.Bookmark
	if err := database.DB.Order("sort_order ASC, created_at ASC").Preload("Category").Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch bookmarks"})
		return
	}
	c.JSON(200, bookmarks)
}

type CreateBookmarkRequest struct {
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required"`
	CategoryID  *uint  `json:"category_id"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

func CreateBookmark(c *gin.Context) {
	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	categoryID := req.CategoryID
	if categoryID == nil {
		var defaultCategory models.Category
		if err := database.DB.Where("is_default = ?", true).First(&defaultCategory).Error; err != nil {
			c.JSON(500, gin.H{"error": "Default category not found"})
			return
		}
		categoryID = &defaultCategory.ID
	}

	var category models.Category
	if err := database.DB.First(&category, *categoryID).Error; err != nil {
		c.JSON(400, gin.H{"error": "Category not found"})
		return
	}

	icon := req.Icon
	if icon == "" {
		fetchedIcon, err := FetchFavicon(req.URL)
		if err == nil {
			icon = fetchedIcon
		}
	}

	bookmark := models.Bookmark{
		CategoryID:  categoryID,
		Title:       req.Title,
		URL:         req.URL,
		Description: req.Description,
		Icon:        icon,
		SortOrder:   req.SortOrder,
		VisitCount:  0,
	}
	if err := database.DB.Create(&bookmark).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create bookmark"})
		return
	}
	c.JSON(201, bookmark)
}

func UpdateBookmark(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var bookmark models.Bookmark
	if err := database.DB.First(&bookmark, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bookmark not found"})
		return
	}

	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if req.CategoryID != nil {
		if err := database.DB.First(&category, *req.CategoryID).Error; err != nil {
			c.JSON(400, gin.H{"error": "Category not found"})
			return
		}
	}

	bookmark.Title = req.Title
	bookmark.URL = req.URL
	if req.CategoryID != nil {
		bookmark.CategoryID = req.CategoryID
	}
	bookmark.Description = req.Description
	bookmark.SortOrder = req.SortOrder
	if req.Icon != "" {
		bookmark.Icon = req.Icon
	} else if bookmark.URL != req.URL {
		fetchedIcon, err := FetchFavicon(req.URL)
		if err == nil {
			bookmark.Icon = fetchedIcon
		}
	}

	if err := database.DB.Save(&bookmark).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update bookmark"})
		return
	}
	c.JSON(200, bookmark)
}

func DeleteBookmark(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var bookmark models.Bookmark
	if err := database.DB.First(&bookmark, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bookmark not found"})
		return
	}

	if err := database.DB.Delete(&bookmark).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete bookmark"})
		return
	}
	c.JSON(200, gin.H{"message": "Bookmark deleted successfully"})
}

func IncrementVisit(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var bookmark models.Bookmark
	if err := database.DB.First(&bookmark, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bookmark not found"})
		return
	}

	if err := database.DB.Model(&bookmark).UpdateColumn("visit_count", bookmark.VisitCount+1).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to increment visit count"})
		return
	}
	c.JSON(200, gin.H{"message": "Visit count incremented"})
}
