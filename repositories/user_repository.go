package repositories

import (
	"backend/models"

	"gorm.io/gorm"
)

// UserRepository - interface untuk operasi pada model User
type UserRepository interface {
	FindByCode(code string) (*models.User, error)
	FindByID(id uint) (*models.User, error)
}

type userRepo struct {
	db *gorm.DB
}

// NewUserRepository - inisialisasi UserRepository baru
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepo{db: db}
}

// FindByCode - mencari user berdasarkan kode unik (code)
func (r *userRepo) FindByCode(code string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("code = ?", code).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByID - mencari user berdasarkan ID
func (r *userRepo) FindByID(id uint) (*models.User, error) {
	var user models.User
	if err := r.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
