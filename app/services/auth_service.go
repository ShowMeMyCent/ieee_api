package services

import (
	"backend/app/models"
	"backend/utils/token"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func LoginCheckAdmin(db *gorm.DB, code string) (string, error) {
	var err error

	adminFromDB := models.User{}
	err = db.Model(&models.User{}).Where("code = ?", code).Take(&adminFromDB).Error
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(adminFromDB.Code), []byte(code))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", err
	}

	token, err := token.GenerateToken(adminFromDB.ID, adminFromDB.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}
