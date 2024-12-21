package routes

import (
	"backend/app/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

func AuthGroup(r *gin.Engine) *gin.RouterGroup {
	authGroup := r.Group("/")
	authGroup.Use(middlewares.AdminCheckMiddleware())
	return authGroup
}

func InitRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Set("db", db)
	})

	r.Use(middlewares.CorsMiddleware())

	// Jika bukan GET , Cek token dulu

	RegisterAuthRoutes(r, db)
	RegisterActivityRoutes(r, db)
	RegisterNewsRoutes(r, db)
	RegisterAchievementRoutes(r, db)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/swagger/doc.json")))
	r.Static("/uploads", "./uploads")
	return r
}
