package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"nav-backend/database"
	"nav-backend/handlers"
	"nav-backend/middleware"
	"nav-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed dist/*
var distFS embed.FS

func main() {
	// 解析命令行参数
	port := flag.String("port", "8081", "Port for web access")
	flag.Parse()

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

	// 获取嵌入的dist文件系统
	embeddedFS, err := fs.Sub(distFS, "dist")
	if err != nil {
		log.Fatal("Failed to create sub filesystem:", err)
	}

	// 提供静态文件服务
	r.GET("/assets/*filepath", func(c *gin.Context) {
		c.FileFromFS(c.Request.URL.Path, http.FS(embeddedFS))
	})

	// 提供index.html
	r.GET("/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		data, err := fs.ReadFile(embeddedFS, "index.html")
		if err != nil {
			c.String(404, "Not found")
			return
		}
		c.Data(200, "text/html; charset=utf-8", data)
	})

	// 处理所有未匹配的路由，用于前端路由
	r.NoRoute(func(c *gin.Context) {
		path := c.Request.URL.Path
		
		// 跳过API路由
		if len(path) >= 4 && path[:4] == "/api" {
			c.JSON(404, gin.H{"error": "Not found"})
			return
		}
		
		// 对于其他路径，提供index.html用于前端路由
		c.Header("Content-Type", "text/html; charset=utf-8")
		data, err := fs.ReadFile(embeddedFS, "index.html")
		if err != nil {
			c.String(404, "Not found")
			return
		}
		c.Data(200, "text/html; charset=utf-8", data)
	})

	// 使用命令行参数或环境变量指定的端口
	portValue := utils.GetEnv("PORT", *port)
	log.Printf("Server starting on port %s...", portValue)
	r.Run(":" + portValue)
}
