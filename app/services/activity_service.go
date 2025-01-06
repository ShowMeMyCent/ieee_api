package services

import (
	"backend/app/models"
	"backend/app/repositories"
)

type ActivityService struct {
	Repo *repositories.ActivityRepository
}

func (as *ActivityService) CreateActivity(activity *models.Activities) error {
	return as.Repo.CreateActivity(activity)
}

func (as *ActivityService) GetAllActivities() ([]models.Activities, error) {
	return as.Repo.GetAllActivities()
}

func (as *ActivityService) GetActivityById(id string) (*models.Activities, error) {
	return as.Repo.GetActivityById(id)
}

func (as *ActivityService) UpdateActivity(activity *models.Activities) error {
	return as.Repo.UpdateActivity(activity)
}

func (as *ActivityService) DeleteActivity(id string) error {
	return as.Repo.DeleteActivity(id)
}
