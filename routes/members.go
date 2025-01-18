package routes

import (
	"backend/app/controllers"
	"backend/app/repositories"
	"backend/app/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterMemberRoutes(r *gin.Engine, db *gorm.DB) {
	memberRepo := &repositories.MembersRepository{DB: db}
	memberService := &services.MemberService{Repo: memberRepo}
	memberController := &controllers.MemberController{Service: memberService}

	r.GET("/members", memberController.GetMembers)
	r.GET("member/:id", memberController.GetMemberById)

	authGroup := AuthGroup(r)
	authGroup.POST("/members", memberController.CreateMember)
	authGroup.PATCH("/member/:id", memberController.UpdateMember)
	authGroup.DELETE("/member/:id")
}
