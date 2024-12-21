package routes

import (
	"backend/app/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterAuthRoutes(r *gin.Engine, db *gorm.DB) {
	r.POST("/login-admin", controllers.LoginAdmin)
}
