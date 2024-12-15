package routes

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/controllers"
	"backend/middlewares"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	paperController := &controllers.PaperController{DB: db}
	activityController := &controllers.ActivitiesController{DB: db}
	achievementController := &controllers.AchievementsController{DB: db}
	newsController := &controllers.NewsController{DB: db}
	r.Use(middlewares.CorsMiddleware())

	// Jika bukan GET , Cek token dulu
	r.Use(func(c *gin.Context) {
		if c.Request.Method != "GET" && c.Request.URL.Path != "/login-admin" {
			middlewares.AdminCheckMiddleware()(c)
		}
		c.Next()
	})

	r.POST("/login-admin", controllers.LoginAdmin)

	r.GET("/papers", paperController.GetAllPapers)
	r.GET("/papers/file/:id", paperController.GetPaperFile)

	//activitiy
	r.GET("/activities", activityController.GetAllActivities)
	r.GET("/activities/:id", activityController.GetActivityById)
	r.GET("/activities/file/:id", activityController.GetGambarActivities)
	//news
	r.GET("/news", newsController.GetAllNews)
	r.GET("/news/:id", newsController.GetNewsById)
	r.GET("/news/file/:id", newsController.GetThumbnailNews)
	r.GET("/news/category/:category", newsController.GetNewsByCategory)

	//achievement
	r.GET("/achievements", achievementController.GetAllAchievement)
	r.GET("/achievements/:id", achievementController.GetAchievementById)
	r.GET("/achievements/foto/:id", achievementController.GetFotoAchievement)
	r.GET("/achievements/category/:category", achievementController.GetAchievementsByCategory)

	//papers
	r.POST("/papers", paperController.UploadPaper)
	r.PUT("/papers/:id", paperController.EditPaper)
	r.DELETE("/papers/:id", paperController.DeletePaper)

	//activities
	r.POST("/activities", activityController.UploadActivity)
	r.DELETE("/activities/:id", activityController.DeleteActivity)
	r.PUT("/activities/:id", activityController.EditActivity)

	//news
	r.POST("/news", newsController.InsertNews)
	r.PUT("/news/:id", newsController.EditNews)
	r.DELETE("/news/:id", newsController.DeleteNews)

	//achievements
	r.POST("/achievements", achievementController.InsertAchievement)
	r.PUT("/achievements/:id", achievementController.EditAchievements)
	r.DELETE("/achievements/:id", achievementController.DeleteAchievements)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	r.Static("/uploads", "./uploads")
	return r
}
