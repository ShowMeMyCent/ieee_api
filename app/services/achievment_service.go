package services

import (
	"backend/app/models"
	"backend/app/repositories"
)

type AchievementService struct {
	Repo *repositories.AchievementRepository
}

func (as *AchievementService) GetAllAchievements() ([]models.Achievement, error) {
	return as.Repo.GetAllAchievements()
}

func (as *AchievementService) GetAchievementById(id string) (*models.Achievement, error) {
	return as.Repo.GetAchievementById(id)
}

func (as *AchievementService) GetAchievementsByCategory(category string) ([]models.Achievement, error) {
	return as.Repo.GetAchievementsByCategory(category)
}

func (as *AchievementService) InsertAchievement(achievement *models.Achievement) error {
	return as.Repo.InsertAchievement(achievement)
}

func (as *AchievementService) UpdateAchievement(achievement *models.Achievement) error {
	return as.Repo.UpdateAchievement(achievement)
}

func (as *AchievementService) DeleteAchievementById(id string) error {
	return as.Repo.DeleteAchievementById(id)
}
