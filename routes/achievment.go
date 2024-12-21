package routes

import (
	"backend/app/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAchievementRoutes(r *gin.Engine, db *gorm.DB) {
	achievementController := &controllers.AchievementsController{DB: db}

	r.GET("/achievements", achievementController.GetAllAchievement)
	r.GET("/achievements/:id", achievementController.GetAchievementById)
	r.GET("/achievements/foto/:id", achievementController.GetFotoAchievement)
	r.GET("/achievements/category/:category", achievementController.GetAchievementsByCategory)

	authGroup := AuthGroup(r)
	authGroup.POST("/achievements", achievementController.InsertAchievement)
	authGroup.PUT("/achievements/:id", achievementController.EditAchievements)
	authGroup.DELETE("/achievements/:id", achievementController.DeleteAchievements)
}
