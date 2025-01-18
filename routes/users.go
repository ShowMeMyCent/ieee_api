package routes

import (
	"backend/app/controllers"
	"backend/app/repositories"
	"backend/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterUserRoutes(r *gin.Engine, db *gorm.DB) {
	userRepo := &repositories.UsersRepository{DB: db}
	userService := &services.UserService{Repo: userRepo}
	userController := &controllers.UserController{Service: userService}

	r.GET("/users", userController.GetUsers)
	r.GET("user/:id", userController.GetUserById)

	authGroup := AuthGroup(r)
	authGroup.POST("/users", userController.CreateUser)
	authGroup.PATCH("/user/:id", userController.UpdateUser)
	authGroup.DELETE("/user/:id", userController.DeleteUser)
}
