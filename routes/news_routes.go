package routes

import (
	"backend/controllers"

	"github.com/gin-gonic/gin"
)

func RegisterNewsRoutes(router *gin.Engine, newsController *controllers.NewsController) {
	news := router.Group("/news")
	{
		news.GET("/", newsController.GetAllNews)
		news.POST("/", newsController.InsertNews)
	}
}
