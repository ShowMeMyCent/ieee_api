package controllers

import (
	"backend/app/models"
	"backend/app/services"
	"backend/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type UserController struct {
	DB      *gorm.DB
	Service *services.UserService
}

type UserResponse struct {
	ID       uint   `gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"-"`
	Role     Role   `json:"role" gorm:"type:enum('admin', 'user', 'it'); default:'user'"`
}

type Role string

func (mc *UserController) GetUsers(ctx *gin.Context) {
	users, err := mc.Service.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}

	var response []UserResponse
	for _, user := range users {
		response = append(response, UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  Role(user.Role),
		})
	}
	ctx.JSON(http.StatusOK, gin.H{"data": response})
}

func (mc *UserController) GetUserById(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := mc.Service.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get user"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"data": user})
}

func (mc *UserController) CreateUser(ctx *gin.Context) {
	// Struktur untuk input JSON
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	validRoles := map[string]bool{"admin": true, "member": true, "tim_it": true}
	if input.Role != "" && !validRoles[input.Role] {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password", "details": err.Error()})
		return
	}

	user := models.User{
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword,
		Role:     models.Role(input.Role),
	}

	validationErrors := utils.ValidateStruct(user)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := mc.Service.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user to database", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (mc *UserController) UpdateUser(ctx *gin.Context) {
	userId := ctx.Param("id")
	user, err := mc.Service.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
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
		user.Name = *input.Name
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Password != nil {
		user.Password = *input.Password
	}

	if input.Role != nil {
		validRoles := map[Role]bool{"admin": true, "user": true, "it": true}
		if !validRoles[*input.Role] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid role"})
			return
		}
		user.Role = models.Role(*input.Role)
	}

	validationErrors := utils.ValidateStruct(user)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := mc.Service.UpdateUser(user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update user to database",
			"details": err.Error(),
		})
		return
	}

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  Role(user.Role),
	}
	ctx.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("User ID: %s updated successfully", userId), "data": response})
}

func (mc *UserController) DeleteUser(ctx *gin.Context) {
	userId := ctx.Param("id")

	_, err := mc.Service.GetUserById(userId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if err := mc.Service.DeleteUser(userId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user", "details": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "User deleted successfully", "data": nil})
}
