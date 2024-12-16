package routes

import (
	"backend/controllers"
	"backend/middlewares"
	"backend/repositories"
	"backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter mengatur semua rute utama aplikasi
func SetupRouter(router *gin.Engine, db *gorm.DB) {
	// Middleware global
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Inisialisasi service dan controller untuk Authentication
	authRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthAdminService(authRepo)
	authController := controllers.NewAuthAdminController(authService)

	// Inisialisasi repository dan service untuk Activities
	activityRepo := repositories.NewActivityRepository(db)       // Inisialisasi repository terlebih dahulu
	activityService := services.NewActivityService(activityRepo) // Gunakan repository untuk service
	activityController := controllers.NewActivitiesController(activityService)

	// Inisialisasi repository dan service untuk News
	newsRepo := repositories.NewNewsRepository(db)
	newsService := services.NewNewsService(newsRepo)
	newsController := controllers.NewNewsController(newsService)

	// Pendaftaran routes modular
	RegisterAuthRoutes(router, authController)
	RegisterActivityRoutes(router, activityController)
	RegisterNewsRoutes(router, newsController)

	// Middleware tambahan (contoh: middleware autentikasi untuk protected routes)
	protected := router.Group("/protected")
	protected.Use(middlewares.AuthMiddleware()) // Hanya untuk rute yang membutuhkan autentikasi
	{
		// Misalnya, kita ingin menambahkan rute yang memerlukan autentikasi
		protected.GET("/secure-data", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"message": "This is protected data!"})
		})
	}
}
