package controllers

import (
	"backend/app/models"
	"backend/app/services"
	"backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivitiesController struct {
	DB      *gorm.DB
	Service *services.ActivityService
}

// UploadActivity adalah fungsi untuk mengupload activity beserta file-nya.
// @Summary Upload an activity with File
// @Description Uploads a Activities along with its file and saves them to the database.
// @Tags Activities
// @Accept multipart/form-data
// @Param title formData string true "Title Activities"
// @Param tanggal formData string true "Tanggal Activities"
// @Param gambar formData file true "File gambar"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.Activities
// @Router /activities [post]
func (ac *ActivitiesController) UploadActivity(ctx *gin.Context) {
	file, err := ctx.FormFile("gambar")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	host := utils.Getenv("ENV_HOST", "localhost")

	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("uploads", fileName)

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	title := ctx.PostForm("title")
	tanggal := ctx.PostForm("tanggal")

	activity := models.Activities{
		Title:    title,
		Tanggal:  tanggal,
		Gambar:   fileName,
		ImageURL: fmt.Sprintf("https://%s/uploads/%s", host, fileName),
	}

	if err := ac.Service.CreateActivity(&activity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save activity to database"})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// GetAllActivities adalah fungsi untuk mendapatkan semua activity dari database.
// @Summary Get All Activities
// @Description Retrieves all Activities from the database.
// @Tags Activities
// @Produce json
// @Success 200 {array} models.Activities
// @Router /activities [get]
func (ac *ActivitiesController) GetAllActivities(ctx *gin.Context) {
	activities, err := ac.Service.GetAllActivities()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get activities from database"})
		return
	}

	ctx.JSON(http.StatusOK, activities)
}

// GetActivityById adalah fungsi untuk mengambil activity berdasarkan ID.
// @Summary Get Activity By ID
// @Description Retrieves Activity data by its ID.
// @Tags Activities
// @Param id path string true "Activity ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /activities/{id} [get]
func (ac *ActivitiesController) GetActivityById(ctx *gin.Context) {
	activityId := ctx.Param("id")

	activity, err := ac.Service.GetActivityById(activityId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// GetGambar adalah fungsi untuk mengambil gambar activity berdasarkan ID.
// @Summary Get Gambar Activity
// @Description Retrieves the image of a activity by its ID.
// @Tags Activities
// @Param id path string true "Activity ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /activities/file/{id} [get]
func (ac *ActivitiesController) GetGambarActivities(ctx *gin.Context) {
	activityID := ctx.Param("id")

	activity, err := ac.Service.GetActivityById(activityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	filePath := filepath.Join("uploads", activity.Gambar)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", activity.Gambar))

	ext := filepath.Ext(activity.Gambar)
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

// EditActivity adalah fungsi untuk mengedit Activity, termasuk kemampuan untuk mengganti file Activity.
// @Summary Edit Activity
// @Description Edits a Activity including the ability to replace its file.
// @Tags Activities
// @Accept multipart/form-data
// @Param id path string true "Activity ID"
// @Param title formData string true "Title Activity"
// @Param tanggal formData string true "Tanggal Activity"
// @Param gambar formData file false "Gambar Activity (optional)"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.Activities
// @Router /activities/{id} [put]
func (ac *ActivitiesController) EditActivity(ctx *gin.Context) {
	activityID := ctx.Param("id")

	activity, err := ac.Service.GetActivityById(activityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	title := ctx.PostForm("title")
	tanggal := ctx.PostForm("tanggal")
	host := utils.Getenv("ENV_HOST", "localhost")

	activity.Title = title
	activity.Tanggal = tanggal

	file, err := ctx.FormFile("gambar")
	if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if file != nil {
		fileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join("uploads", fileName)

		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		oldFilePath := filepath.Join("uploads", activity.Gambar)
		if err := os.Remove(oldFilePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old file"})
		}

		activity.Gambar = fileName
		activity.ImageURL = fmt.Sprintf("https://%s/uploads/%s", host, fileName)
	}

	if err := ac.Service.UpdateActivity(activity); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated activity to database"})
		return
	}

	ctx.JSON(http.StatusOK, activity)
}

// DeleteActivity adalah fungsi untuk menghapus Activity dan gambar-nya dari database.
// @Summary Delete Activity
// @Description Deletes a Activity and its gambar from the database and storage.
// @Tags Activities
// @Param id path string true "Activity ID"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {string} string "Activity deleted successfully"
// @Router /activities/{id} [delete]
func (ac *ActivitiesController) DeleteActivity(ctx *gin.Context) {
	activityID := ctx.Param("id")

	activity, err := ac.Service.GetActivityById(activityID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Activity not found"})
		return
	}

	filePath := filepath.Join("uploads", activity.Gambar)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		if err := ac.Service.DeleteActivity(activityID); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity from database"})
			return
		}
		ctx.JSON(http.StatusOK, "Activity deleted successfully")
		return
	}

	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	if err := ac.Service.DeleteActivity(activityID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete activity from database"})
		return
	}

	ctx.JSON(http.StatusOK, "Activity and file deleted successfully")
}
