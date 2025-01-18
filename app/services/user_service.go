package services

import (
	"backend/app/models"
	"backend/app/repositories"
)

type UserService struct {
	Repo *repositories.UsersRepository
}

func (ms *UserService) GetUsers() ([]models.User, error) {
	return ms.Repo.GetUsers()
}

func (ms *UserService) GetUserById(id string) (*models.User, error) {
	return ms.Repo.GetUsersById(id)
}

func (ms *UserService) CreateUser(user *models.User) error {
	return ms.Repo.CreateUsers(user)
}

func (ms *UserService) UpdateUser(user *models.User) error {
	return ms.Repo.UpdateUser(user)
}

func (ms *UserService) DeleteUser(id string) error {
	return ms.Repo.DeleteUser(id)
}
