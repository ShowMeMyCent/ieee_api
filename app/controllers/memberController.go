package controllers

import (
	"backend/app/models"
	"backend/app/services"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type MemberController struct {
	DB      *gorm.DB
	Service *services.MemberService
}

type MembersResponse struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"-"`
	Role     Role   `json:"role" gorm:"type:enum('admin', 'member', 'it'); default:'member'"`
}

type Role string

func (mc *MemberController) GetMembers(ctx *gin.Context) {
	members, err := mc.Service.GetMembers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get members"})
		return
	}

	var response []MembersResponse
	for _, member := range members {
		response = append(response, MembersResponse{
			ID:    member.ID,
			Name:  member.Name,
			Email: member.Email,
			Role:  Role(member.Role),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (mc *MemberController) GetMemberById(ctx *gin.Context) {
	memberId := ctx.Param("id")
	member, err := mc.Service.GetMemberById(memberId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get member"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": member})
}

func (mc *MemberController) CreateMember(ctx *gin.Context) {
	// Struktur untuk input JSON
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Role     Role   `json:"role" gorm:"type:enum('admin', 'member', 'it'); default:'member'"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	//validRoles := map[Role]bool{"admin": true, "member": true, "it": true}
	//if !validRoles[input.Role] {
	//	ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
	//	return
	//}

	member := models.Members{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
		Role:     models.Role(input.Role),
	}

	validationErrors := utils.ValidateStruct(member)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := mc.Service.CreateMember(&member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save member to database", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Member created successfully"})
}

func (mc *MemberController) UpdateMember(ctx *gin.Context) {
	memberId := ctx.Param("id")
	member, err := mc.Service.GetMemberById(memberId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Member not found"})
		return
	}

	var input struct {
		Name     *string `json:"name" binding:"omitempty"`
		Email    *string `json:"email" binding:"omitempty,email"`
		Password *string `json:"password" binding:"omitempty"`
		Role     *Role   `json:"role" binding:"omitempty"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	if input.Name != nil {
		member.Name = *input.Name
	}
	if input.Email != nil {
		member.Email = *input.Email
	}
	if input.Password != nil {
		member.Password = *input.Password
	}

	if input.Role != nil {
		validRoles := map[Role]bool{"admin": true, "member": true, "it": true}
		if !validRoles[*input.Role] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
		member.Role = models.Role(*input.Role)
	}

	validationErrors := utils.ValidateStruct(member)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := mc.Service.UpdateMember(member); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update member to database",
			"details": err.Error(),
		})
		return
	}

	response := MembersResponse{
		ID:    member.ID,
		Name:  member.Name,
		Email: member.Email,
		Role:  Role(member.Role),
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Member %s updated successfully", memberId), "data": response})
}
