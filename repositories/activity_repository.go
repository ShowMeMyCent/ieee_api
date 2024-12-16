package repositories

import (
	"backend/models"
	"errors"
	"gorm.io/gorm"
)

type ActivityRepository interface {
	Create(activity *models.Activities) error
	FindAll() ([]models.Activities, error)
	FindByID(id string) (models.Activities, error)
	Update(id, title, tanggal, fileName string) (models.Activities, error)
	Delete(id string) error
}

type activityRepo struct {
	db *gorm.DB
}

func NewActivityRepository(db *gorm.DB) ActivityRepository {
	return &activityRepo{db: db}
}

func (r *activityRepo) Create(activity *models.Activities) error {
	return r.db.Create(activity).Error
}

func (r *activityRepo) FindAll() ([]models.Activities, error) {
	var activities []models.Activities
	err := r.db.Find(&activities).Error
	return activities, err
}

func (r *activityRepo) FindByID(id string) (models.Activities, error) {
	var activity models.Activities
	err := r.db.Where("id = ?", id).First(&activity).Error
	return activity, err
}

func (r *activityRepo) Update(id, title, tanggal, fileName string) (models.Activities, error) {
	var activity models.Activities
	if err := r.db.Where("id = ?", id).First(&activity).Error; err != nil {
		return activity, errors.New("activity not found")
	}

	activity.Title = title
	activity.Date = tanggal
	if fileName != "" {
		activity.Image = fileName
		activity.ImageURL = "uploads/" + fileName
	}

	if err := r.db.Save(&activity).Error; err != nil {
		return activity, err
	}
	return activity, nil
}

func (r *activityRepo) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Activities{}).Error
}
