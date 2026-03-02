package handlers

import (
	"nav-backend/database"
	"nav-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	userID := c.GetUint("user_id")

	var categories []models.Category
	if err := database.DB.Where("user_id = ?", userID).Order("sort_order ASC, created_at ASC").Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch categories"})
		return
	}

	c.JSON(200, categories)
}

type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func CreateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
		SortOrder:   0,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create category"})
		return
	}

	c.JSON(201, category)
}

func UpdateCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(404, gin.H{"error": "Category not found"})
		return
	}

	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category.Name = req.Name
	category.Description = req.Description
	category.Icon = req.Icon

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update category"})
		return
	}

	c.JSON(200, category)
}

func DeleteCategory(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := database.DB.Where("id = ? AND user_id = ?", id, userID).First(&category).Error; err != nil {
		c.JSON(404, gin.H{"error": "Category not found"})
		return
	}

	tx := database.DB.Begin()
	
	if err := tx.Delete(&category).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to delete category"})
		return
	}

	if err := tx.Where("category_id = ?", id).Delete(&models.Bookmark{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to delete associated bookmarks"})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Category deleted successfully"})
}
