package routes

import (
	"backend/app/controllers/authentication"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRouts(r *gin.Engine, db *gorm.DB) {
	authController := &authentication.AuthController{DB: db}

	authGroup := r.Group("/auth")
	{
		authGroup.POST("/login", authController.Login)
		authGroup.POST("/register", authController.Register)
	}
}

//OLD
//func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
//	r.POST("/login-admin", controllers.LoginAdmin)
//}
