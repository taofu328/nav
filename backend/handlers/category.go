package handlers

import (
	"nav-backend/database"
	"nav-backend/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category
	if err := database.DB.Order("is_default DESC, sort_order ASC, created_at ASC").Find(&categories).Error; err != nil {
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
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	category := models.Category{
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
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
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
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := database.DB.First(&category, id).Error; err != nil {
		c.JSON(404, gin.H{"error": "Category not found"})
		return
	}

	if category.IsDefault {
		c.JSON(400, gin.H{"error": "Cannot delete default category"})
		return
	}

	var defaultCategory models.Category
	if err := database.DB.Where("is_default = ?", true).First(&defaultCategory).Error; err != nil {
		c.JSON(500, gin.H{"error": "Default category not found"})
		return
	}

	tx := database.DB.Begin()
	
	if err := tx.Model(&models.Bookmark{}).Where("category_id = ?", id).Update("category_id", defaultCategory.ID).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to migrate bookmarks"})
		return
	}

	if err := tx.Delete(&category).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "Failed to delete category"})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "Category deleted successfully", "migrated_to": defaultCategory.ID})
}
