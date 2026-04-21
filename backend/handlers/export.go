package handlers

import (
	"nav-backend/database"
	"nav-backend/logger"
	"nav-backend/models"

	"github.com/gin-gonic/gin"
)

// 导入导出数据结构
type ExportData struct {
	Categories []ExportCategory `json:"categories"`
	Bookmarks  []ExportBookmark `json:"bookmarks"`
}

// 分类导出数据结构
type ExportCategory struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	IsDefault   bool   `json:"is_default"`
}

// 书签导出数据结构
type ExportBookmark struct {
	ID          uint   `json:"id"`
	CategoryID  *uint  `json:"category_id"`
	Title       string `json:"title"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
	VisitCount  int    `json:"visit_count"`
}

// 导入请求结构
type ImportRequest struct {
	Data     ExportData `json:"data"`
	Conflict string     `json:"conflict"` // 冲突处理策略: "overwrite", "skip", "merge"
}

// 导出数据
func ExportDataHandler(c *gin.Context) {
	// 导出分类
	var categories []models.Category
	if err := database.DB.Order("is_default DESC, sort_order ASC").Find(&categories).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to export categories"})
		return
	}

	// 导出书签
	var bookmarks []models.Bookmark
	if err := database.DB.Order("sort_order ASC").Preload("Category").Find(&bookmarks).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to export bookmarks"})
		return
	}

	// 转换为导出格式
	exportCategories := make([]ExportCategory, len(categories))
	for i, category := range categories {
		exportCategories[i] = ExportCategory{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			Icon:        category.Icon,
			SortOrder:   category.SortOrder,
			IsDefault:   category.IsDefault,
		}
	}

	exportBookmarks := make([]ExportBookmark, len(bookmarks))
	for i, bookmark := range bookmarks {
		exportBookmarks[i] = ExportBookmark{
			ID:          bookmark.ID,
			CategoryID:  bookmark.CategoryID,
			Title:       bookmark.Title,
			URL:         bookmark.URL,
			Description: bookmark.Description,
			Icon:        bookmark.Icon,
			SortOrder:   bookmark.SortOrder,
			VisitCount:  bookmark.VisitCount,
		}
	}

	exportData := ExportData{
		Categories: exportCategories,
		Bookmarks:  exportBookmarks,
	}

	logger.Info("ExportData: exported %d categories and %d bookmarks", len(exportCategories), len(exportBookmarks))

	c.JSON(200, exportData)
}

// 导入数据
func ImportDataHandler(c *gin.Context) {
	var req ImportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request format"})
		return
	}

	tx := database.DB.Begin()

	// 处理分类导入
	categoryMap := make(map[uint]*uint)       // 原始ID -> 新ID
	categoryNameMap := make(map[string]*uint) // 分类名称 -> ID

	// 先导入所有分类
	for _, category := range req.Data.Categories {
		// 检查分类是否已存在
		var existingCategory models.Category
		err := tx.Where("name = ?", category.Name).First(&existingCategory).Error

		if err == nil {
			// 分类已存在，根据冲突策略处理
			switch req.Conflict {
			case "overwrite":
				// 覆盖现有分类
				existingCategory.Description = category.Description
				existingCategory.Icon = category.Icon
				existingCategory.SortOrder = category.SortOrder
				existingCategory.IsDefault = category.IsDefault
				if err := tx.Save(&existingCategory).Error; err != nil {
					tx.Rollback()
					c.JSON(500, gin.H{"error": "Failed to update category"})
					return
				}
				categoryMap[category.ID] = &existingCategory.ID
				categoryNameMap[category.Name] = &existingCategory.ID
				logger.Info("ImportData: Updated category '%s'", category.Name)
			case "skip":
				// 跳过现有分类
				categoryMap[category.ID] = &existingCategory.ID
				categoryNameMap[category.Name] = &existingCategory.ID
				logger.Info("ImportData: Skipped existing category '%s'", category.Name)
			case "merge":
				// 合并分类（使用现有分类）
				categoryMap[category.ID] = &existingCategory.ID
				categoryNameMap[category.Name] = &existingCategory.ID
				logger.Info("ImportData: Merged with existing category '%s'", category.Name)
			default:
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid conflict strategy"})
				return
			}
		} else {
			// 分类不存在，创建新分类
			newCategory := models.Category{
				Name:        category.Name,
				Description: category.Description,
				Icon:        category.Icon,
				SortOrder:   category.SortOrder,
				IsDefault:   category.IsDefault,
			}
			if err := tx.Create(&newCategory).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to create category"})
				return
			}
			categoryMap[category.ID] = &newCategory.ID
			categoryNameMap[category.Name] = &newCategory.ID
			logger.Info("ImportData: Created new category '%s' with ID %d", category.Name, newCategory.ID)
		}
	}

	// 处理书签导入
	bookmarkCount := 0
	for _, bookmark := range req.Data.Bookmarks {
		// 检查书签是否已存在（通过URL判断）
		var existingBookmark models.Bookmark
		err := tx.Where("url = ?", bookmark.URL).First(&existingBookmark).Error

		if err == nil {
			// 书签已存在，根据冲突策略处理
			switch req.Conflict {
			case "overwrite":
				// 覆盖现有书签
				existingBookmark.Title = bookmark.Title
				existingBookmark.Description = bookmark.Description
				existingBookmark.Icon = bookmark.Icon
				existingBookmark.SortOrder = bookmark.SortOrder
				existingBookmark.VisitCount = bookmark.VisitCount
				// 更新分类ID
				if bookmark.CategoryID != nil {
					if newCategoryID, exists := categoryMap[*bookmark.CategoryID]; exists {
						existingBookmark.CategoryID = newCategoryID
					}
				}
				if err := tx.Save(&existingBookmark).Error; err != nil {
					tx.Rollback()
					c.JSON(500, gin.H{"error": "Failed to update bookmark"})
					return
				}
				bookmarkCount++
				logger.Info("ImportData: Updated bookmark '%s'", bookmark.Title)
			case "skip":
				// 跳过现有书签
				logger.Info("ImportData: Skipped existing bookmark '%s'", bookmark.Title)
			case "merge":
				// 合并书签（使用现有书签）
				logger.Info("ImportData: Merged with existing bookmark '%s'", bookmark.Title)
			default:
				tx.Rollback()
				c.JSON(400, gin.H{"error": "Invalid conflict strategy"})
				return
			}
		} else {
			// 书签不存在，创建新书签
			newBookmark := models.Bookmark{
				Title:       bookmark.Title,
				URL:         bookmark.URL,
				Description: bookmark.Description,
				Icon:        bookmark.Icon,
				SortOrder:   bookmark.SortOrder,
				VisitCount:  bookmark.VisitCount,
			}
			// 设置分类ID
			if bookmark.CategoryID != nil {
				if newCategoryID, exists := categoryMap[*bookmark.CategoryID]; exists {
					newBookmark.CategoryID = newCategoryID
				}
			}
			if err := tx.Create(&newBookmark).Error; err != nil {
				tx.Rollback()
				c.JSON(500, gin.H{"error": "Failed to create bookmark"})
				return
			}
			bookmarkCount++
			logger.Info("ImportData: Created new bookmark '%s'", bookmark.Title)
		}
	}

	tx.Commit()
	logger.Info("ImportData: Imported %d categories and %d bookmarks", len(req.Data.Categories), bookmarkCount)

	c.JSON(200, gin.H{
		"message":    "Data imported successfully",
		"categories": len(req.Data.Categories),
		"bookmarks":  bookmarkCount,
		"conflict":   req.Conflict,
	})
}
