package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ActivitiesController menangani request terkait activities
type ActivitiesController struct {
	Service services.ActivityService
}

// NewActivitiesController membuat instance ActivitiesController
func NewActivitiesController(service services.ActivityService) *ActivitiesController {
	return &ActivitiesController{Service: service}
}

// UploadActivity mengunggah aktivitas beserta file gambar
func (ac *ActivitiesController) UploadActivity(ctx *gin.Context) {
	// Ambil file dari form-data
	file, err := ctx.FormFile("gambar")
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "File gambar is required")
		return
	}

	// Validasi file gambar
	if err := utils.ValidateFile(file); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Simpan file ke folder uploads
	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("uploads", fileName)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to save file")
		return
	}

	// Ambil input form lainnya
	title := ctx.PostForm("title")
	tanggal := ctx.PostForm("tanggal")
	host := utils.Getenv("ENV_HOST", "localhost")

	// Buat object activity
	activity := models.Activities{
		Title:    title,
		Date:     tanggal,
		Image:    fileName,
		ImageURL: fmt.Sprintf("https://%s/uploads/%s", host, fileName),
	}

	// Simpan ke database
	if err := ac.Service.CreateActivity(&activity); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to save activity to database")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity uploaded successfully", "data": activity})
}

// GetAllActivities mengambil semua aktivitas
func (ac *ActivitiesController) GetAllActivities(ctx *gin.Context) {
	activities, err := ac.Service.GetAllActivities()
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to get activities")
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

// GetActivityByID mengambil aktivitas berdasarkan ID
func (ac *ActivitiesController) GetActivityByID(ctx *gin.Context) {
	id := ctx.Param("id")
	activity, err := ac.Service.GetActivityByID(id)
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Activity not found")
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// EditActivity mengupdate data aktivitas
func (ac *ActivitiesController) EditActivity(ctx *gin.Context) {
	id := ctx.Param("id")

	// Ambil data form
	title := ctx.PostForm("title")
	tanggal := ctx.PostForm("tanggal")

	// Ambil file (jika ada)
	file, err := ctx.FormFile("gambar")
	var fileName string
	if err == nil { // File ditemukan
		// Validasi file
		if err := utils.ValidateFile(file); err != nil {
			utils.HandleError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Simpan file baru
		fileName = uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join("uploads", fileName)
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Failed to save file")
			return
		}
	}

	// Panggil service untuk update
	updatedActivity, err := ac.Service.UpdateActivity(id, title, tanggal, fileName)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to update activity")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity updated successfully", "data": updatedActivity})
}

// DeleteActivity menghapus aktivitas dan gambarnya
func (ac *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	id := ctx.Param("id")

	// Panggil service untuk hapus
	if err := ac.Service.DeleteActivity(id); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Failed to delete activity")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
