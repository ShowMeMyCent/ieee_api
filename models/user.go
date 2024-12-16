package models

import (
	"backend/utils/token"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User - model untuk user
type User struct {
	ID   uint   `gorm:"primaryKey;autoIncrement:true"`
	Role string `json:"role"`
	Code string `json:"code"`
}

// GenerateJWT - generate token JWT untuk user
func (u *User) GenerateJWT() (string, error) {
	if u.ID == 0 || u.Role == "" {
		return "", errors.New("invalid user data for token generation")
	}
	return token.GenerateToken(u.ID, u.Role)
}

// LoginCheckAdmin - verifikasi kode admin
func (u *User) LoginCheckAdmin(db *gorm.DB) (string, error) {
	var admin User
	if err := db.Where("code = ?", u.Code).First(&admin).Error; err != nil {
		return "", errors.New("admin not found")
	}

	// Validasi kode dengan bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(admin.Code), []byte(u.Code)); err != nil {
		return "", errors.New("invalid code")
	}

	// Generate token JWT
	return admin.GenerateJWT()
}
