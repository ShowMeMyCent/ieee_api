package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// AchievementController menangani request terkait achievements
type AchievementController struct {
	Service services.AchievementService
}

// NewAchievementController membuat instance baru dari AchievementController
func NewAchievementController(service services.AchievementService) *AchievementController {
	return &AchievementController{Service: service}
}

// GetAllAchievements mengambil semua achievement
func (ac *AchievementController) GetAllAchievements(ctx *gin.Context) {
	achievements, err := ac.Service.GetAllAchievements()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve achievements"})
		return
	}
	ctx.JSON(http.StatusOK, achievements)
}

// GetAchievementByID mengambil achievement berdasarkan ID
func (ac *AchievementController) GetAchievementByID(ctx *gin.Context) {
	id := ctx.Param("id")

	achievement, err := ac.Service.GetAchievementByID(id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}
	ctx.JSON(http.StatusOK, achievement)
}

// CreateAchievement menambahkan achievement baru
func (ac *AchievementController) CreateAchievement(ctx *gin.Context) {
	var achievement models.Achievement

	// Bind JSON ke struct Achievement
	if err := ctx.ShouldBindJSON(&achievement); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Simpan achievement
	if err := ac.Service.CreateAchievement(&achievement); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create achievement"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Achievement created successfully", "data": achievement})
}

// UploadAchievementImage mengunggah foto untuk achievement
func (ac *AchievementController) UploadAchievementImage(ctx *gin.Context) {
	// Ambil file dari form-data
	file, err := ctx.FormFile("foto")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	// Validasi file
	if err := utils.ValidateFile(file); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate file name dan simpan
	fileName := uuid.New().String() + ".png"
	filePath := "uploads/" + fileName

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Dapatkan ID dari URL parameter
	id := ctx.Param("id")

	// Update path gambar di database
	if err := ac.Service.UpdateAchievementImage(id, fileName, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update image path"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully", "path": filePath})
}

// UpdateAchievement mengupdate achievement berdasarkan ID
func (ac *AchievementController) UpdateAchievement(ctx *gin.Context) {
	id := ctx.Param("id")
	var updatedData models.Achievement

	// Bind input
	if err := ctx.ShouldBindJSON(&updatedData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Update achievement
	if err := ac.Service.UpdateAchievement(id, &updatedData); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievement"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Achievement updated successfully", "data": updatedData})
}

// DeleteAchievement menghapus achievement berdasarkan ID
func (ac *AchievementController) DeleteAchievement(ctx *gin.Context) {
	id := ctx.Param("id")

	// Hapus achievement
	if err := ac.Service.DeleteAchievement(id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievement"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Achievement deleted successfully"})
}
