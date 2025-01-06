package routes

import (
	"backend/app/controllers"
	"backend/app/repositories"
	"backend/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterActivityRoutes(r *gin.Engine, db *gorm.DB) {
	activityRepo := &repositories.ActivityRepository{DB: db}
	activityService := &services.ActivityService{Repo: activityRepo}
	activityController := &controllers.ActivitiesController{Service: activityService}

	r.GET("/activities", activityController.GetAllActivities)
	r.GET("/activities/:id", activityController.GetActivityById)
	r.GET("/activities/file/:id", activityController.GetGambarActivities)

	authGroup := AuthGroup(r)
	authGroup.POST("/activities", activityController.UploadActivity)
	authGroup.PUT("/activities/:id", activityController.EditActivity)
	authGroup.DELETE("/activities/:id", activityController.DeleteActivity)
}
