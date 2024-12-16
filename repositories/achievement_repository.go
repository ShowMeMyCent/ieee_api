package repositories

import (
	"backend/models"
	"errors"

	"gorm.io/gorm"
)

// AchievementRepository mendefinisikan metode untuk interaksi database
type AchievementRepository interface {
	FindAll() ([]models.Achievement, error)
	FindByID(id string) (models.Achievement, error)
	Create(achievement *models.Achievement) error
	Update(id string, updatedData *models.Achievement) error
	UpdateImagePath(id string, filePath string) error
	Delete(id string) error
}

type achievementRepo struct {
	db *gorm.DB
}

// NewAchievementRepository membuat instance baru dari AchievementRepository
func NewAchievementRepository(db *gorm.DB) AchievementRepository {
	return &achievementRepo{db: db}
}

// FindAll mengambil semua data achievement dari database
func (repo *achievementRepo) FindAll() ([]models.Achievement, error) {
	var achievements []models.Achievement
	if err := repo.db.Find(&achievements).Error; err != nil {
		return nil, err
	}
	return achievements, nil
}

// FindByID mencari achievement berdasarkan ID
func (repo *achievementRepo) FindByID(id string) (models.Achievement, error) {
	var achievement models.Achievement
	if err := repo.db.Where("id = ?", id).First(&achievement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return achievement, errors.New("achievement not found")
		}
		return achievement, err
	}
	return achievement, nil
}

// Create menambahkan data achievement ke database
func (repo *achievementRepo) Create(achievement *models.Achievement) error {
	if err := repo.db.Create(achievement).Error; err != nil {
		return err
	}
	return nil
}

// Update memperbarui data achievement berdasarkan ID
func (repo *achievementRepo) Update(id string, updatedData *models.Achievement) error {
	var achievement models.Achievement
	if err := repo.db.Where("id = ?", id).First(&achievement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("achievement not found")
		}
		return err
	}

	// Update field yang diubah
	if err := repo.db.Model(&achievement).Updates(updatedData).Error; err != nil {
		return err
	}
	return nil
}

// UpdateImagePath memperbarui path file gambar dari achievement
func (repo *achievementRepo) UpdateImagePath(id string, filePath string) error {
	var achievement models.Achievement
	if err := repo.db.Where("id = ?", id).First(&achievement).Error; err != nil {
		return err
	}

	// Update field foto dan image_url
	achievement.Foto = filePath
	achievement.ImageURL = filePath

	if err := repo.db.Save(&achievement).Error; err != nil {
		return err
	}
	return nil
}

// Delete menghapus data achievement dari database
func (repo *achievementRepo) Delete(id string) error {
	var achievement models.Achievement
	if err := repo.db.Where("id = ?", id).First(&achievement).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("achievement not found")
		}
		return err
	}

	// Hapus dari database
	if err := repo.db.Delete(&achievement).Error; err != nil {
		return err
	}
	return nil
}
