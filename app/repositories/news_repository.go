package repositories

import (
	"backend/app/models"
	"gorm.io/gorm"
)

type NewsRepository struct {
	DB *gorm.DB
}

func (nr *NewsRepository) GetAllNews() ([]models.News, error) {
	var news []models.News
	err := nr.DB.Find(&news).Error
	return news, err
}

func (nr *NewsRepository) GetNewsById(id string) (*models.News, error) {
	var news models.News
	err := nr.DB.Where("id = ?", id).First(&news).Error
	return &news, err
}

func (nr *NewsRepository) GetNewsByCategory(category string) ([]models.News, error) {
	var news []models.News
	err := nr.DB.Where("kategori = ?", category).Find(&news).Error
	return news, err
}

func (nr *NewsRepository) CreateNews(news *models.News) error {
	return nr.DB.Create(news).Error
}

func (nr *NewsRepository) UpdateNews(news *models.News) error {
	return nr.DB.Save(news).Error
}

func (nr *NewsRepository) DeleteNewsById(id string) error {
	return nr.DB.Where("id = ?", id).Delete(&models.News{}).Error
}
