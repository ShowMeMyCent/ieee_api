package services

import (
	"backend/models"
	"backend/repositories"
)

type ActivityService interface {
	CreateActivity(activity *models.Activities) error
	GetAllActivities() ([]models.Activities, error)
	GetActivityByID(id string) (models.Activities, error)
	UpdateActivity(id, title, tanggal, fileName string) (models.Activities, error)
	DeleteActivity(id string) error
}

type activityService struct {
	repo repositories.ActivityRepository
}

func NewActivityService(repo repositories.ActivityRepository) ActivityService {
	return &activityService{repo: repo}
}

func (s *activityService) CreateActivity(activity *models.Activities) error {
	return s.repo.Create(activity)
}

func (s *activityService) GetAllActivities() ([]models.Activities, error) {
	return s.repo.FindAll()
}

func (s *activityService) GetActivityByID(id string) (models.Activities, error) {
	return s.repo.FindByID(id)
}

func (s *activityService) UpdateActivity(id, title, tanggal, fileName string) (models.Activities, error) {
	return s.repo.Update(id, title, tanggal, fileName)
}

func (s *activityService) DeleteActivity(id string) error {
	return s.repo.Delete(id)
}
