package repositories

import (
	"backend/app/models"
	"gorm.io/gorm"
)

type MembersRepository struct {
	DB *gorm.DB
}

func (mr *MembersRepository) GetMembers() ([]models.Members, error) {
	var members []models.Members
	if err := mr.DB.Find(&members).Error; err != nil {
		return nil, err
	}
	return members, nil
}

func (mr *MembersRepository) GetMembersById(id string) (*models.Members, error) {
	var member models.Members
	err := mr.DB.Where("id = ?", id).First(&member).Error
	return &member, err
}

func (mr *MembersRepository) CreateMembers(members *models.Members) error {
	return mr.DB.Create(members).Error
}

func (mr *MembersRepository) UpdateMembers(members *models.Members) error {
	return mr.DB.Save(members).Error
}
