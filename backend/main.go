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

	r.POST("/api/auth/register", handlers.Register)
	r.POST("/api/auth/login", handlers.Login)

	r.GET("/api/public/categories", handlers.GetPublicCategories)
	r.GET("/api/public/bookmarks", handlers.GetPublicBookmarks)

	r.GET("/api/icons", handlers.GetFavicon)
	r.GET("/api/icons/:filename", handlers.ServeIcon)

	r.POST("/api/admin/login", handlers.AdminLogin)

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

		api.GET("/export", handlers.ExportData)
		api.POST("/import", handlers.ImportData)

		api.DELETE("/admin/clear-all", handlers.ClearAllData)
	api.POST("/admin/update-icons", handlers.UpdateAllIcons)
	api.POST("/icons/upload", handlers.UploadIcon)
	api.DELETE("/icons", handlers.DeleteIcon)
	}

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Server starting on port %s...", port)
	r.Run(":" + port)
}
