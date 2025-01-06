package controllers

import (
	"backend/app/models"
	"backend/app/services"
	"backend/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AchievementsController struct {
	DB      *gorm.DB
	Service *services.AchievementService
}

// GetAllAchievements adalah fungsi untuk mendapatkan semua achievements dari database.
// @Summary Get All Achievement
// @Description Retrieves all achievements from the database.
// @Tags Achievement
// @Produce json
// @Success 200 {array} models.Achievement
// @Router /achievements [get]
func (ac *AchievementsController) GetAllAchievement(ctx *gin.Context) {
	var achievements []models.Achievement

	// Retrieve all achievementsfrom the database
	if err := ac.DB.Find(&achievements).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get achievements from database"})
		return
	}

	// Response with achievements
	ctx.JSON(http.StatusOK, achievements)
}

// GetAchievementById adalah fungsi untuk mengambil achievements berdasarkan ID.
// @Summary Get achievement By ID
// @Description Retrieves achievement data by its ID.
// @Tags Achievement
// @Param id path string true "achievement ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /achievements/{id} [get]
func (ac *AchievementsController) GetAchievementById(ctx *gin.Context) {
	achievementId := ctx.Param("id")
	achievement, err := ac.Service.GetAchievementById(achievementId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}
	ctx.JSON(http.StatusOK, achievement)
}

// GetFotoAchievement adalah fungsi untuk mengambil Foto Achievement berdasarkan ID.
// @Summary Get Foto Achievement
// @Description Retrieves the image of an Achievement by its ID.
// @Tags Achievement
// @Param id path string true "Achievement ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /achievements/foto/{id} [get]
func (ac *AchievementsController) GetFotoAchievement(ctx *gin.Context) {
	// Get achievements image ID from URL path parameter, pakenya param klo mau diubah ke yg laen sok
	achievementsID := ctx.Param("id")

	// Retrieve Thumbnail from the database by its ID
	var achievements models.Achievement
	if err := ac.DB.Where("id = ?", achievementsID).First(&achievements).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	// Define the file path
	filePath := filepath.Join("uploads", achievements.Foto)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set the headers for the file transfer and return the file
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", achievements.Foto))

	// switch case buat content type
	ext := filepath.Ext(achievements.Foto)
	switch ext {
	case ".png":
		ctx.Header("Content-Type", "image/png")
	case ".jpg", ".jpeg":
		ctx.Header("Content-Type", "image/jpeg")
	case ".gif":
		ctx.Header("Content-Type", "image/gif")
	case ".pdf":
		ctx.Header("Content-Type", "application/pdf")
	default:
		ctx.Header("Content-Type", "application/octet-stream")
	}

	ctx.File(filePath)
}

// GetAchievementsByCategory adalah fungsi untuk mengambil Achievements berdasarkan Category.
// @Summary Get All Achievements By Category
// @Description Retrieves Achievements data by its Category.
// @Tags Achievement
// @Param category path string true "Achievements Category"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /achievements/category/{category} [get]
func (ac *AchievementsController) GetAchievementsByCategory(ctx *gin.Context) {
	category := ctx.Param("category")
	achievements, err := ac.Service.GetAchievementsByCategory(category)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievements not found"})
		return
	}
	ctx.JSON(http.StatusOK, achievements)
}

// InsertAchievement adalah fungsi untuk membuat post Achievements terbaru.
// @Summary Insert a new Achievement
// @Description Insert a Achievements and saves them to the database.
// @Tags Achievement
// @Accept multipart/form-data
// @Param nama formData string true "Nama peraih achievement"
// @Param pencapaian formData string true "Pencapaian yang diraih"
// @Param link formData string true "link ke Achievementnya"
// @Param kategori formData string true "Kategori Achievement"
// @Param foto formData file true "Foto peraih Achievement"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.Achievement
// @Router /achievements [post]
func (ac *AchievementsController) InsertAchievement(ctx *gin.Context) {
	file, err := ctx.FormFile("foto")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	host := utils.Getenv("ENV_HOST", "localhost")

	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("uploads", fileName)

	var achievement models.Achievement
	achievement.Nama = ctx.PostForm("nama")
	achievement.Pencapaian = ctx.PostForm("pencapaian")
	achievement.Link = ctx.PostForm("link")
	achievement.Kategori = ctx.PostForm("kategori")
	achievement.Foto = fileName
	achievement.ImageURL = fmt.Sprintf("https://%s/uploads/%s", host, fileName)

	validationErrors := utils.ValidateStruct(achievement)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if err := ac.Service.InsertAchievement(&achievement); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save achievement to database"})
		return
	}

	ctx.JSON(http.StatusOK, achievement)
}

// EditAchievements adalah fungsi untuk mengedit Achievements
// @Summary Edit Achievements
// @Description Edits a Achievements by its ID
// @Tags Achievement
// @Accept multipart/form-data
// @Param id path string true "Achievement ID"
// @Param nama formData string true "Nama peraih achievement"
// @Param pencapaian formData string true "Pencapaian yang diraih"
// @Param link formData string true "link ke Achievementnya"
// @Param kategori formData string true "Kategori Achievement"
// @Param foto formData file true "Foto peraih Achievement"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.Achievement
// @Router /achievements/{id} [put]
func (ac *AchievementsController) EditAchievement(ctx *gin.Context) {
	achievementId := ctx.Param("id")

	var achievement *models.Achievement
	achievement, err := ac.Service.GetAchievementById(achievementId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	achievement.Nama = ctx.PostForm("nama")
	achievement.Pencapaian = ctx.PostForm("pencapaian")
	achievement.Link = ctx.PostForm("link")
	achievement.Kategori = ctx.PostForm("kategori")

	validationErrors := utils.ValidateStruct(achievement)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	file, err := ctx.FormFile("foto")
	if err == nil {
		fileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join("uploads", fileName)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		oldFilePath := filepath.Join("uploads", achievement.Foto)
		if err := os.Remove(oldFilePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old file"})
			return
		}

		achievement.Foto = fileName
		achievement.ImageURL = fmt.Sprintf("http://localhost:8080/uploads/%s", fileName)
	}

	if err := ac.Service.UpdateAchievement(achievement); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement in database"})
		return
	}

	ctx.JSON(http.StatusOK, achievement)
}

// DeleteAchievements adalah fungsi untuk menghapus Achievements dan gambar-nya dari database.
// @Summary Delete Achievements
// @Description Delete a Achievements and its thumbnail from the database and storage.
// @Tags Achievement
// @Param id path string true "Achievements ID"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {string} string "Achievements deleted successfully"
// @Router /achievements/{id} [delete]
func (ac *AchievementsController) DeleteAchievement(ctx *gin.Context) {
	achievementId := ctx.Param("id")

	achievement, err := ac.Service.GetAchievementById(achievementId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	filePath := filepath.Join("uploads", achievement.Foto)
	if _, err := os.Stat(filePath); err == nil {
		if err := os.Remove(filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
			return
		}
	}

	if err := ac.Service.DeleteAchievementById(achievementId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievement from database"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
