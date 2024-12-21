package routes

import (
	"backend/app/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterActivityRoutes(r *gin.Engine, db *gorm.DB) {
	activityController := &controllers.ActivitiesController{DB: db}

	r.GET("/activities", activityController.GetAllActivities)
	r.GET("/activities/:id", activityController.GetActivityById)
	r.GET("/activities/file/:id", activityController.GetGambarActivities)

	authGroup := AuthGroup(r)
	authGroup.POST("/activities", activityController.UploadActivity)
	authGroup.PUT("/activities/:id", activityController.EditActivity)
	authGroup.DELETE("/activities/:id", activityController.DeleteActivity)
}
