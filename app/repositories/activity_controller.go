package repositories

import (
	"backend/app/models"
	"gorm.io/gorm"
)

type ActivityRepository struct {
	DB *gorm.DB
}

func (ar *ActivityRepository) CreateActivity(activity *models.Activities) error {
	return ar.DB.Create(activity).Error
}

func (ar *ActivityRepository) GetAllActivities() ([]models.Activities, error) {
	var activities []models.Activities
	err := ar.DB.Find(&activities).Error
	return activities, err
}

func (ar *ActivityRepository) GetActivityById(id string) (*models.Activities, error) {
	var activity models.Activities
	err := ar.DB.Where("id = ?", id).First(&activity).Error
	return &activity, err
}

func (ar *ActivityRepository) UpdateActivity(activity *models.Activities) error {
	return ar.DB.Save(activity).Error
}

func (ar *ActivityRepository) DeleteActivity(id string) error {
	return ar.DB.Where("id = ?", id).Delete(&models.Activities{}).Error
}
