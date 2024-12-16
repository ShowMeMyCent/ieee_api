package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterAchievementRoutes(router *gin.Engine, controller *controllers.AchievementController) {
	achievements := router.Group("/achievements")
	{
		achievements.GET("/", controller.GetAllAchievements)
		achievements.GET("/:id", controller.GetAchievementByID)
		achievements.POST("/", controller.CreateAchievement)
		achievements.POST("/:id/upload-image", controller.UploadAchievementImage)
		achievements.PUT("/:id", controller.UpdateAchievement)
		achievements.DELETE("/:id", controller.DeleteAchievement)
	}
}
