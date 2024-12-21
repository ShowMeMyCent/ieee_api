package routes

import (
	"backend/app/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNewsRoutes(r *gin.Engine, db *gorm.DB) {
	newsController := &controllers.NewsController{DB: db}

	r.GET("/news", newsController.GetAllNews)
	r.GET("/news/:id", newsController.GetNewsById)
	r.GET("/news/file/:id", newsController.GetThumbnailNews)
	r.GET("/news/category/:category", newsController.GetNewsByCategory)

	authGroup := AuthGroup(r)
	authGroup.POST("/news", newsController.InsertNews)
	authGroup.PUT("/news/:id", newsController.EditNews)
	authGroup.DELETE("/news/:id", newsController.DeleteNews)
}
