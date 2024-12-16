package services

import (
	"backend/models"
	"backend/repositories"
)

type AchievementService interface {
	GetAllAchievements() ([]models.Achievement, error)
	GetAchievementByID(id string) (models.Achievement, error)
	CreateAchievement(achievement *models.Achievement) error
	UpdateAchievement(id string, updatedData *models.Achievement) error
	UpdateAchievementImage(id, fileName, filePath string) error
	DeleteAchievement(id string) error
}

type achievementService struct {
	repo repositories.AchievementRepository
}

func NewAchievementService(repo repositories.AchievementRepository) AchievementService {
	return &achievementService{repo: repo}
}

func (s *achievementService) GetAllAchievements() ([]models.Achievement, error) {
	return s.repo.FindAll()
}

func (s *achievementService) GetAchievementByID(id string) (models.Achievement, error) {
	return s.repo.FindByID(id)
}

func (s *achievementService) CreateAchievement(achievement *models.Achievement) error {
	return s.repo.Create(achievement)
}

func (s *achievementService) UpdateAchievement(id string, updatedData *models.Achievement) error {
	return s.repo.Update(id, updatedData)
}

func (s *achievementService) UpdateAchievementImage(id, fileName, filePath string) error {
	return s.repo.UpdateImagePath(id, filePath)
}

func (s *achievementService) DeleteAchievement(id string) error {
	return s.repo.Delete(id)
}
