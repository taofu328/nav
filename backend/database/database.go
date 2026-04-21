package database

import (
	"nav-backend/config"
	"nav-backend/logger"
	"nav-backend/models"
	"nav-backend/utils"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = config.DBPath
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Bookmark{},
		&models.SiteSetting{},
	)
	if err != nil {
		logger.Fatal("Failed to migrate database: %v", err)
	}

	// 检查是否需要创建默认管理员账户
	createDefaultAdminIfNeeded()

	logger.Info("Database initialized successfully")
}

func createDefaultAdminIfNeeded() {
	var userCount int64
	if err := DB.Model(&models.User{}).Count(&userCount).Error; err != nil {
		logger.Error("Failed to check user count: %v", err)
		return
	}

	if userCount > 0 {
		logger.Info("Users already exist, skipping default admin creation")
		return
	}

	logger.Info("No users found, creating default admin account...")

	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		logger.Error("Failed to hash admin password: %v", err)
		return
	}

	admin := models.User{
		Username: "admin",
		Email:    "admin@nav.local",
		Password: hashedPassword,
	}

	if err := DB.Create(&admin).Error; err != nil {
		logger.Error("Failed to create admin account: %v", err)
	} else {
		logger.Info("Default admin account created successfully")
	}
}
