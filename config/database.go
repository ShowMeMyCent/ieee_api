package config

import (
	"backend/models"
	"backend/utils"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

// DatabaseConfig stores the database configuration
type DatabaseConfig struct {
	Username string
	Password string
	Database string
	Host     string
}

// GetDatabaseConfig retrieves the database configuration from environment variables
func GetDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		Username: utils.Getenv("DB_USERNAME", ""),
		Password: utils.Getenv("DB_PASSWORD", ""),
		Database: utils.Getenv("DB_DATABASE", ""),
		Host:     utils.Getenv("DB_HOST", ""),
	}
}

// ConnectDatabase establishes a connection to the database
func ConnectDatabase() *gorm.DB {
	// Get the database configuration
	config := GetDatabaseConfig()

	// Create the Data Source Name (DSN) string for the connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Username,
		config.Password,
		config.Host,
		config.Database,
	)

	// Open a connection to the MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	// Automatically migrate the database schema
	AutoMigrate(db)

	return db
}

// AutoMigrate performs automatic migration of database schema
func AutoMigrate(db *gorm.DB) {
	// Run migrations for the models
	err := db.AutoMigrate(
		&models.Achievement{},
		&models.Activities{},
		&models.News{},
		&models.User{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
	}
}
