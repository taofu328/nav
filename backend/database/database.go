package database

import (
	"database/sql"
	"log"
	"nav-backend/config"
	"nav-backend/models"
	"nav-backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	_ "modernc.org/sqlite"
)

var DB *gorm.DB

func InitDB() {
	dbPath := config.DBPath
	if dbPath == "" {
		dbPath = config.DBPath
	}

	var err error
	sqlDB, err := sql.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	DB, err = gorm.Open(sqlite.Dialector{
		Conn: sqlDB,
	}, &gorm.Config{
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
	// 检查是否已经存在admin用户
	var adminUser models.User
	result := DB.Where("username = ?", "admin").First(&adminUser)

	// 打印查询结果
	if result.Error == nil {
		log.Println("Admin user already exists:", adminUser.Username)
		return
	}

	// 尝试通过邮箱查找用户
	result = DB.Where("email = ?", "admin@nav.local").First(&adminUser)

	// 创建默认管理员账户
	log.Println("Creating default admin account...")

	// 创建默认管理员账户
	hashedPassword, err := utils.HashPassword("admin")
	if err != nil {
		log.Println("Failed to hash admin password:", err)
		return
	}

	if result.Error == nil {
		// 更新现有用户的用户名和密码
		adminUser.Username = "admin"
		adminUser.Password = hashedPassword
		if err := DB.Save(&adminUser).Error; err != nil {
			log.Println("Failed to update admin account:", err)
		} else {
			log.Println("Default admin account updated successfully")
		}
	} else {
		// 创建新用户
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
