package main

import (
	"flag"
	"log"
	"nav-backend/database"
	"nav-backend/handlers"
	"nav-backend/middleware"
	"nav-backend/utils"
	"os"

	"github.com/gin-gonic/gin"
)

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

	// 提供静态文件服务，处理所有未匹配的路由
	r.NoRoute(func(c *gin.Context) {
		// 尝试提供静态文件
		filePath := "./dist" + c.Request.URL.Path
		if _, err := os.Stat(filePath); err == nil {
			// 文件存在，直接提供
			c.File(filePath)
		} else {
			// 文件不存在，提供index.html，用于前端路由
			c.File("./dist/index.html")
		}
	})

	// 使用命令行参数或环境变量指定的端口
	portValue := utils.GetEnv("PORT", *port)
	log.Printf("Server starting on port %s...", portValue)
	r.Run(":" + portValue)
}
