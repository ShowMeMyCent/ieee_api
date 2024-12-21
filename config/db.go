package config

import (
	"backend/app/models"
	"backend/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func ConnectDatabase() *gorm.DB {
	// EDIT SESUAI DB LOKAL
	username := utils.Getenv("DB_USERNAME", "root")
	password := utils.Getenv("DB_PASSWORD", "12345")
	database := utils.Getenv("DB_DATABASE", "yourdb")
	host := utils.Getenv("DB_HOST", "127.0.0.1")

	dsn := fmt.Sprintf("%s:%s@(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username,
		password,
		host,
		database,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	db.AutoMigrate(&models.Achievement{}, &models.Activities{}, &models.News{}, &models.User{})

	return db
}
