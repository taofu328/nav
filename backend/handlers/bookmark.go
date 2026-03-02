package handlers

import (
	"nav-backend/database"
	"nav-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetBookmarks(c *gin.Context) {
	userID := c.GetUint("user_id")

	var bookmarks []models.Bookmark
	if err := database.DB.Where("user_id = ?", userID).Order("sort_order ASC, created_at ASC").Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch bookmarks"})
		return
	}

	c.JSON(200, bookmarks)
}

type CreateBookmarkRequest struct {
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func CreateBookmark(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := database.DB.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
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
		UserID:      userID,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		URL:         req.URL,
		Description: req.Description,
		Icon:        icon,
		SortOrder:   0,
		VisitCount:  0,
	}

	if err := database.DB.Create(&bookmark).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create bookmark"})
		return
	}

	c.JSON(201, bookmark)
}

func UpdateBookmark(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var bookmark models.Bookmark
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&bookmark).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bookmark not found"})
		return
	}

	var req CreateBookmarkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := database.DB.Where("id = ? AND user_id = ?", req.CategoryID, userID).First(&category).Error; err != nil {
		c.JSON(400, gin.H{"error": "Category not found"})
		return
	}

	bookmark.Title = req.Title
	bookmark.URL = req.URL
	bookmark.CategoryID = req.CategoryID
	bookmark.Description = req.Description
	
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
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var bookmark models.Bookmark
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&bookmark).Error; err != nil {
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
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var bookmark models.Bookmark
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&bookmark).Error; err != nil {
		c.JSON(404, gin.H{"error": "Bookmark not found"})
		return
	}

	if err := database.DB.Model(&bookmark).UpdateColumn("visit_count", bookmark.VisitCount+1).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update visit count"})
		return
	}

	c.JSON(200, gin.H{"message": "Visit count updated"})
}
