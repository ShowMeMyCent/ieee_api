package repositories

import (
	"backend/app/models"
	"gorm.io/gorm"
)

type AchievementRepository struct {
	DB *gorm.DB
}

func (ar *AchievementRepository) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := ar.DB.Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (ar *AchievementRepository) GetAchievementById(id string) (*models.Achievement, error) {
	var achievement models.Achievement
	if err := ar.DB.Where("id = ?", id).First(&achievement).Error; err != nil {
		return nil, err
	}
	return &achievement, nil
}

func (ar *AchievementRepository) GetAchievementsByCategory(category string) ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := ar.DB.Where("kategori = ?", category).Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

func (ar *AchievementRepository) InsertAchievement(achievement *models.Achievement) error {
	return ar.DB.Create(achievement).Error
}

func (ar *AchievementRepository) UpdateAchievement(achievement *models.Achievement) error {
	return ar.DB.Save(achievement).Error
}

func (ar *AchievementRepository) DeleteAchievementById(id string) error {
	return ar.DB.Where("id = ?", id).Delete(&models.Achievement{}).Error
}
