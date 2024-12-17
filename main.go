package main

import (
	"backend/config"
	"backend/controllers"
	"backend/repositories"
	"backend/routes"
	"backend/services"
	"backend/utils"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to the database
	db := config.ConnectDatabase()
	if db == nil {
		log.Fatal("Failed to connect to the database")
	}

	// Initialize Gin router
	router := gin.Default()

	// Initialize repositories and services
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthAdminService(userRepo)
	authController := controllers.NewAuthAdminController(authService)

	activityService := services.NewActivityService(db)
	activityController := controllers.NewActivitiesController(activityService)

	newsRepo := repositories.NewNewsRepository(db)
	newsService := services.NewNewsService(newsRepo)
	newsController := controllers.NewNewsController(newsService)

	// Setup routes
	routes.RegisterAuthRoutes(router, authController)
	routes.RegisterActivityRoutes(router, activityController)
	routes.RegisterNewsRoutes(router, newsController)

	// Run the server
	host := utils.Getenv("ENV_HOST", "localhost")
	port := utils.Getenv("ENV_PORT", "8080")
	serverAddr := fmt.Sprintf("%s:%s", host, port)

	log.Printf("Server is running on %s", serverAddr)
	if err := router.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
