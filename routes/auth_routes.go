package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

func RegisterAuthRoutes(router *gin.Engine, authController *controllers.AuthAdminController) {
	auth := router.Group("/auth")
	{
		auth.POST("/login-admin", authController.LoginAdmin)
	}
}
