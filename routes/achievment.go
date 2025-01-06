package routes

import (
	"backend/app/controllers"
	"backend/app/repositories"
	"backend/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAchievementRoutes(r *gin.Engine, db *gorm.DB) {
	achievementRepo := &repositories.AchievementRepository{DB: db}
	achievementService := &services.AchievementService{Repo: achievementRepo}
	achievementController := &controllers.AchievementsController{Service: achievementService}

	r.GET("/achievements", achievementController.GetAllAchievement)
	r.GET("/achievements/:id", achievementController.GetAchievementById)
	r.GET("/achievements/foto/:id", achievementController.GetFotoAchievement)
	r.GET("/achievements/category/:category", achievementController.GetAchievementsByCategory)

	authGroup := AuthGroup(r)
	authGroup.POST("/achievements", achievementController.InsertAchievement)
	authGroup.PUT("/achievements/:id", achievementController.EditAchievement)
	authGroup.DELETE("/achievements/:id", achievementController.DeleteAchievement)
}
