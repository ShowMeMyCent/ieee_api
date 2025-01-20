package repositories

import (
	"backend/app/models"
	"gorm.io/gorm"
)

type UsersRepository struct {
	DB *gorm.DB
}

func (mr *UsersRepository) GetUsers() ([]models.User, error) {
	var Users []models.User
	if err := mr.DB.Find(&Users).Error; err != nil {
		return nil, err
	}
	return Users, nil
}

func (mr *UsersRepository) GetUsersById(id string) (*models.User, error) {
	var user models.User
	err := mr.DB.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (mr *UsersRepository) CreateUsers(Users *models.User) error {
	return mr.DB.Create(Users).Error
}

func (mr *UsersRepository) UpdateUser(Users *models.User) error {
	return mr.DB.Save(Users).Error
}

func (mr *UsersRepository) DeleteUser(id string) error {
	return mr.DB.Where("id = ?", id).Delete(&models.User{}).Error
}
