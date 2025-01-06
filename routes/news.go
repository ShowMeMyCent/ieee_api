package routes

import (
	"backend/app/controllers"
	"backend/app/repositories"
	"backend/app/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterNewsRoutes(r *gin.Engine, db *gorm.DB) {
	newsRepo := &repositories.NewsRepository{DB: db}
	newsService := &services.NewsService{Repo: newsRepo}
	newsController := &controllers.NewsController{Service: newsService}

	r.GET("/news", newsController.GetAllNews)
	r.GET("/news/:id", newsController.GetNewsById)
	r.GET("/news/file/:id", newsController.GetThumbnailNews)
	r.GET("/news/category/:category", newsController.GetNewsByCategory)

	authGroup := AuthGroup(r)
	authGroup.POST("/news", newsController.InsertNews)
	authGroup.PUT("/news/:id", newsController.EditNews)
	authGroup.DELETE("/news/:id", newsController.DeleteNews)
}
