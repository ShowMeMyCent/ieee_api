package config

import (
	"backend/app/models"
	"gorm.io/gorm"
	"log"
)

// AutoMigrate performs automatic migration of database schema
func AutoMigrate(db *gorm.DB) {
	// Run migrations for the models
	err := db.AutoMigrate(
		&models.Achievement{},
		&models.Activities{},
		&models.News{},
		&models.User{},
		&models.Members{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
	}
}
