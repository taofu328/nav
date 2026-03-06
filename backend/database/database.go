package database

import (
	"log"
	"nav-backend/config"
	"nav-backend/models"
	"nav-backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = config.DBPath
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Bookmark{},
		&models.SiteSetting{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// 检查是否需要创建默认管理员账户
	createDefaultAdminIfNeeded()

	log.Println("Database initialized successfully")
}

func createDefaultAdminIfNeeded() {
	// 检查 users 表是否为空
	var userCount int64
	if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		log.Println("Failed to check user count:", err)
		return
	}

	// 只有当没有用户时才创建默认管理员账户
	if userCount == 0 {
		log.Println("Creating default admin account...")
		
		// 创建默认管理员账户
		hashedPassword, err := utils.HashPassword("admin")
		if err != nil {
			log.Println("Failed to hash admin password:", err)
			return
		}

		admin := models.User{
			Username: "admin",
			Email:    "admin@nav.local",
			Password: hashedPassword,
		}

		if err := DB.Create(&admin).Error; err != nil {
			log.Println("Failed to create admin account:", err)
		} else {
			log.Println("Default admin account created successfully")
		}
	}
}
