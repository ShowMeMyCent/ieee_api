package authentication

import (
	"backend/app/models"
	"backend/utils"
	token2 "backend/utils/token"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

type AuthController struct {
	DB *gorm.DB
}

type UserLoginResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type DataResponse struct {
	Token string            `json:"token"`
	User  UserLoginResponse `json:"user"`
}

type Response struct {
	Status string       `json:"status"`
	Data   DataResponse `json:"data"`
}

// Login handles user login
func (ac *AuthController) Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// Input Validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find User by Email
	var user models.User
	if err := ac.DB.Where("email = ?", strings.ToLower(input.Email)).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	// Validate Password
	if err := utils.CheckPassword(user.Password, input.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials password"})
		return
	}

	// Generate Token (placeholder function, implement as needed)
	token, err := token2.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	// Respond with token and user data
	c.JSON(http.StatusOK, Response{
		Status: "success",
		Data: DataResponse{
			Token: token,
			User: UserLoginResponse{
				ID:    user.ID,
				Name:  user.Name,
				Email: user.Email,
				Role:  string(user.Role),
			},
		},
	})

}

// Register handles user registration
func (ac *AuthController) Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
		Role     string `json:"role"`
	}

	// Input Validation
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate Role
	role := models.Role(strings.ToLower(input.Role))
	if input.Role != "" && !role.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
		return
	}

	// Hash Password
	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Create User
	user := models.User{
		Name:     input.Name,
		Email:    strings.ToLower(input.Email),
		Password: hashedPassword, // Hash disimpan langsung di sini
		Role:     role,
	}

	// Simpan ke database
	if err := ac.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user created successfully"})
}

//OLD

//type KodeInput struct {
//	Kode string `json:"kode" binding:"required"`
//}
//
//// LoginAdmin godoc
//// @Summary Login as as admin.
//// @Description Logging in to get jwt token to access admin or user api by roles.
//// @Tags Auth
//// @Param Body body KodeInput true "the body to login a admin"
//// @Produce json
//// @Success 200 {object} map[string]interface{}
//// @Router /login-admin [post]
//func LoginAdmin(c *gin.Context) {
//	db := c.MustGet("db").(*gorm.DB)
//	var input KodeInput
//
//	if err := c.ShouldBindJSON(&input); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//
//	token, err := services.LoginCheckAdmin(db, input.Kode)
//
//	if err != nil {
//		fmt.Println(err)
//		c.JSON(http.StatusBadRequest, gin.H{"error": "code is incorrect."})
//		return
//	}
//
//	admin := map[string]string{
//		"kode": input.Kode,
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "login success", "user": admin, "token": token})
//}
