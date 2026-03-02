package handlers

import (
	"nav-backend/database"
	"nav-backend/models"
	"nav-backend/utils"

	"github.com/gin-gonic/gin"
)

type AdminLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdminLogin(c *gin.Context) {
	var req AdminLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var adminUser models.User
	if err := database.DB.Where("username = ?", req.Username).First(&adminUser).Error; err != nil {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	if !utils.CheckPassword(req.Password, adminUser.Password) {
		c.JSON(401, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(adminUser.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"user": gin.H{
			"id":       adminUser.ID,
			"username": adminUser.Username,
			"email":    adminUser.Email,
		},
	})
}

func ClearAllData(c *gin.Context) {
	tx := database.DB.Begin()

	if err := tx.Exec("DELETE FROM bookmarks").Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "清空数据失败"})
		return
	}

	if err := tx.Exec("DELETE FROM categories").Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "清空数据失败"})
		return
	}

	tx.Commit()
	c.JSON(200, gin.H{"message": "数据已清空"})
}
