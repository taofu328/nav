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

type UpdateUserRequest struct {
	Username    string `json:"username" binding:"required"`
	NewPassword string `json:"new_password"`
}

type UpdateSettingsRequest struct {
	SiteTitle string `json:"site_title"`
	SiteLogo  string `json:"site_logo"`
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

func UpdateUserInfo(c *gin.Context) {
	userID := c.GetUint("user_id")
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误"})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(404, gin.H{"error": "用户不存在"})
		return
	}

	user.Username = req.Username
	
	if req.NewPassword != "" {
		hashedPassword, err := utils.HashPassword(req.NewPassword)
		if err != nil {
			c.JSON(500, gin.H{"error": "密码加密失败"})
			return
		}
		user.Password = hashedPassword
	}

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "更新用户信息失败"})
		return
	}

	c.JSON(200, gin.H{
		"message": "更新成功",
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func GetSiteSettings(c *gin.Context) {
	// 从数据库中获取 site_setting 记录
	var siteSetting models.SiteSetting
	result := database.DB.First(&siteSetting, 1)

	// 如果记录不存在，返回默认设置
	if result.Error != nil {
		c.JSON(200, gin.H{
			"settings": gin.H{
				"site_title": "Van Nav",
				"site_logo":  "",
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"settings": gin.H{
			"site_title": siteSetting.Title,
			"site_logo":  siteSetting.LogoURL,
		},
	})
}

func UpdateSiteSettings(c *gin.Context) {
	var req UpdateSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "请求参数错误"})
		return
	}

	// 检查数据库中是否已有 site_setting 记录
	var siteSetting models.SiteSetting
	result := database.DB.First(&siteSetting, 1)

	// 如果记录不存在，创建一个新记录
	if result.Error != nil {
		siteSetting = models.SiteSetting{
			ID:      1,
			Title:   req.SiteTitle,
			LogoURL: req.SiteLogo,
		}
		if err := database.DB.Create(&siteSetting).Error; err != nil {
			c.JSON(500, gin.H{"error": "保存设置失败"})
			return
		}
	} else {
		// 如果记录存在，更新它
		siteSetting.Title = req.SiteTitle
		siteSetting.LogoURL = req.SiteLogo
		if err := database.DB.Save(&siteSetting).Error; err != nil {
			c.JSON(500, gin.H{"error": "保存设置失败"})
			return
		}
	}

	c.JSON(200, gin.H{
		"message": "设置已保存",
		"settings": gin.H{
			"site_title": siteSetting.Title,
			"site_logo":  siteSetting.LogoURL,
		},
	})
}
