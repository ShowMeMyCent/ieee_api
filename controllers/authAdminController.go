package controllers

import (
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// KodeInput - struct untuk input kode login
type KodeInput struct {
	Kode string `json:"kode" binding:"required,min=5,max=30"`
}

// AuthAdminController - controller untuk autentikasi admin
type AuthAdminController struct {
	AuthService services.AuthAdminService
}

// NewAuthAdminController - inisialisasi controller baru
func NewAuthAdminController(authService services.AuthAdminService) *AuthAdminController {
	return &AuthAdminController{AuthService: authService}
}

// LoginAdmin - endpoint untuk autentikasi admin
// @Summary Login as admin.
// @Description Logs in an admin and returns a JWT token.
// @Tags Auth
// @Accept json
// @Param Body KodeInput true "Admin code for login"
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /auth/login-admin [post]
func (ac *AuthAdminController) LoginAdmin(ctx *gin.Context) {
	var input KodeInput

	// Validasi input
	if err := ctx.ShouldBindJSON(&input); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Invalid input format")
		return
	}

	// Panggil service untuk autentikasi
	token, err := ac.AuthService.LoginAdmin(input.Kode)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	// Respons sukses
	ctx.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}
