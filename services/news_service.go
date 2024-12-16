package services

import (
	"backend/models"
	"backend/repositories"
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
)

type NewsService interface {
	GetAllNews(ctx context.Context) ([]models.News, error)
	InsertNews(ctx context.Context, news models.News, file *multipart.FileHeader, ginContext *gin.Context) (models.News, error)
}

type newsService struct {
	repo repositories.NewsRepository
}

func NewNewsService(repo repositories.NewsRepository) NewsService {
	return &newsService{repo: repo}
}

// GetAllNews mengambil seluruh berita
func (s *newsService) GetAllNews(ctx context.Context) ([]models.News, error) {
	return s.repo.FindAll(ctx)
}

// InsertNews menyimpan berita dan thumbnail ke database
func (s *newsService) InsertNews(ctx context.Context, news models.News, file *multipart.FileHeader, ginContext *gin.Context) (models.News, error) {
	// Simpan file thumbnail dan dapatkan nama file
	fileName, err := s.repo.SaveFile(file, ginContext) // Menggunakan ginContext
	if err != nil {
		return models.News{}, errors.New("failed to save file")
	}

	// Update properti gambar pada berita
	news.Thumbnail = fileName
	news.ImageURL = fmt.Sprintf("%s/uploads/%s", ginContext.Request.Host, fileName)

	// Simpan berita ke database
	if err := s.repo.Save(ctx, &news); err != nil {
		return models.News{}, errors.New("failed to save news")
	}

	return news, nil
}
