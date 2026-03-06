package main

import (
	"log"
	"nav-backend/database"
	"nav-backend/handlers"
	"nav-backend/middleware"
	"nav-backend/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()

	r := gin.Default()

	r.Use(middleware.CORS())

	r.POST("/api/auth/login", handlers.Login)

	r.GET("/api/public/categories", handlers.GetPublicCategories)
	r.GET("/api/public/bookmarks", handlers.GetPublicBookmarks)

	r.GET("/api/icons", handlers.GetFavicon)
	r.GET("/api/icons/:filename", handlers.ServeIcon)

	r.POST("/api/admin/login", handlers.AdminLogin)
	// 读取网站设置不需要身份验证
	r.GET("/api/admin/settings", handlers.GetSiteSettings)

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/categories", handlers.GetCategories)
		api.POST("/categories", handlers.CreateCategory)
		api.PUT("/categories/:id", handlers.UpdateCategory)
		api.DELETE("/categories/:id", handlers.DeleteCategory)

		api.GET("/bookmarks", handlers.GetBookmarks)
		api.POST("/bookmarks", handlers.CreateBookmark)
		api.PUT("/bookmarks/:id", handlers.UpdateBookmark)
		api.DELETE("/bookmarks/:id", handlers.DeleteBookmark)
		api.POST("/bookmarks/:id/visit", handlers.IncrementVisit)

		api.GET("/export", handlers.ExportDataHandler)
		api.POST("/import", handlers.ImportDataHandler)

		api.DELETE("/admin/clear-all", handlers.ClearAllData)
		api.POST("/admin/update-icons", handlers.UpdateAllIcons)
		api.POST("/icons/upload", handlers.UploadIcon)
		api.DELETE("/icons", handlers.DeleteIcon)
		api.PUT("/admin/user", handlers.UpdateUserInfo)
		// 修改网站设置需要身份验证
		api.PUT("/admin/settings", handlers.UpdateSiteSettings)
	}

	port := utils.GetEnv("PORT", "8081")
	log.Printf("Server starting on port %s...", port)
	r.Run(":" + port)
}
