package services

import (
	"backend/app/models"
	"backend/app/repositories"
)

type NewsService struct {
	Repo *repositories.NewsRepository
}

func (ns *NewsService) GetAllNews() ([]models.News, error) {
	return ns.Repo.GetAllNews()
}

func (ns *NewsService) GetNewsById(id string) (*models.News, error) {
	return ns.Repo.GetNewsById(id)
}

func (ns *NewsService) GetNewsByCategory(category string) ([]models.News, error) {
	return ns.Repo.GetNewsByCategory(category)
}

func (ns *NewsService) CreateNews(news *models.News) error {
	return ns.Repo.CreateNews(news)
}

func (ns *NewsService) UpdateNews(news *models.News) error {
	return ns.Repo.UpdateNews(news)
}

func (ns *NewsService) DeleteNewsById(id string) error {
	return ns.Repo.DeleteNewsById(id)
}
