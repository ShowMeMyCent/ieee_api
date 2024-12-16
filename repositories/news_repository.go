package repositories

import (
	"backend/models"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// NewsRepository - interface untuk operasi pada model News
type NewsRepository interface {
	FindAll(ctx context.Context) ([]models.News, error)
	Save(ctx context.Context, news *models.News) error
	SaveFile(file *multipart.FileHeader, ginContext *gin.Context) (string, error) // Menggunakan ginContext
}

type newsRepository struct {
	db *gorm.DB
}

func NewNewsRepository(db *gorm.DB) NewsRepository {
	return &newsRepository{db: db}
}

func (r *newsRepository) FindAll(ctx context.Context) ([]models.News, error) {
	var news []models.News
	if err := r.db.Find(&news).Error; err != nil {
		return nil, err
	}
	return news, nil
}

func (r *newsRepository) Save(ctx context.Context, news *models.News) error {
	if err := r.db.Create(news).Error; err != nil {
		return err
	}
	return nil
}

// SaveFile menyimpan file gambar dan mengembalikan nama file
func (r *newsRepository) SaveFile(file *multipart.FileHeader, ginContext *gin.Context) (string, error) {
	// Tentukan direktori tempat file akan disimpan
	dir := "uploads"
	// Membuat folder jika belum ada
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return "", err
	}

	// Menghasilkan nama file unik
	fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), file.Filename)
	filePath := filepath.Join(dir, fileName) // Menggunakan filePath untuk menentukan lokasi penyimpanan

	// Menggunakan ginContext untuk menyimpan file
	if err := ginContext.SaveUploadedFile(file, filePath); err != nil {
		return "", err
	}

	return fileName, nil
}
