package controllers

import (
	"backend/models"
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
	DB *gorm.DB
}

// GetAllAchievements adalah fungsi untuk mendapatkan semua achievements dari database.
// @Summary Get All Achievement
// @Description Retrieves all achievements from the database.
// @Tags Achievement
// @Produce json
// @Success 200 {array} models.Achievement
// @Router /achievements [get]
func (ac *AchievementsController) GetAllAchievement(ctx *gin.Context) {
	var achievements[]models.Achievement

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
	// Get achievement ID from URL path parameter
	achievementsId := ctx.Param("id")

	// Retrieve achievement from the database by its ID
	var achievements models.Achievement
	if err := ac.DB.Where("id = ?", achievementsId).First(&achievements).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "achievements not found"})
		return
	}

	// Response with achievements
	ctx.JSON(http.StatusOK, achievements)
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
	// Get Achievements category from URL path parameter
	achievementsCategory := ctx.Param("category")

	// Retrieve achievements from the database by its ID
	var achievements []models.Achievement
	if err := ac.DB.Where("kategori = ?", achievementsCategory).Find(&achievements).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	// Response with achievements
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
    // Get the file from the form data
    file, err := ctx.FormFile("foto")
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
	host := utils.Getenv("ENV_HOST", "localhost")

    // Define the path where the file will be saved, pake UUID, untuk skg taro di uploads dlu
    fileName := uuid.New().String() + filepath.Ext(file.Filename)
    filePath := filepath.Join("uploads", fileName)

    // Validate input data using struct tags
    var achievements models.Achievement
    achievements.Nama = ctx.PostForm("nama")
    achievements.Pencapaian = ctx.PostForm("pencapaian")
    achievements.Link = ctx.PostForm("link")
    achievements.Kategori = ctx.PostForm("kategori")
    achievements.Foto = fileName
    achievements.ImageURL = fmt.Sprintf("https://%s/uploads/%s", host, fileName)

    // Validate the achievements struct
    validationErrors := utils.ValidateStruct(achievements)
    if len(validationErrors) > 0 {
        // Return validation errors without saving the file
        ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
        return
    }

    // Save the file to the defined path
    if err := ctx.SaveUploadedFile(file, filePath); err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
        return
    }

    // Save achievement to database
    if err := ac.DB.Create(&achievements).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save achievement to database"})
        return
    }

    ctx.JSON(http.StatusOK, achievements)
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
func (ac *AchievementsController) EditAchievements(ctx *gin.Context) {
    // Get achievement ID from URL path parameter
    achievementsID := ctx.Param("id")

    // Retrieve achievement from the database by its ID
    var achievements models.Achievement
    if err := ac.DB.Where("id = ?", achievementsID).First(&achievements).Error; err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{"error": "achievements not found"})
        return
    }

    // Update achievements fields from form data
    achievements.Nama = ctx.PostForm("nama")
    achievements.Pencapaian = ctx.PostForm("pencapaian")
    achievements.Link = ctx.PostForm("link")
    achievements.Kategori = ctx.PostForm("kategori")

    // Validate the achievement struct
    validationErrors := utils.ValidateStruct(achievements)
    if len(validationErrors) > 0 {
        ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
        return
    }

    // Update the file if a new one is uploaded
    file, err := ctx.FormFile("foto")
    if err == nil {
        // Define the path where the new file will be saved
        fileName := uuid.New().String() + filepath.Ext(file.Filename)
        filePath := filepath.Join("uploads", fileName)

        // Save the new file to the defined path
        if err := ctx.SaveUploadedFile(file, filePath); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
            return
        }

        // Delete the old file
        oldFilePath := filepath.Join("uploads", achievements.Foto)
        if err := os.Remove(oldFilePath); err != nil {
            ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete old file"})
            return
        }

        // Update the Foto field with the new file name and update ImageURL
        achievements.Foto = fileName
        achievements.ImageURL = fmt.Sprintf("http://localhost:8080/uploads/%s", fileName)
    }

    // Save updated achievement to the database
    if err := ac.DB.Save(&achievements).Error; err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update achievements in database"})
        return
    }

    ctx.JSON(http.StatusOK, achievements)
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
func (ac *AchievementsController) DeleteAchievements(ctx *gin.Context) {
	achievementId := ctx.Param("id")

	// Cari achievements dari id
	var achievements models.Achievement
	if err := ac.DB.Where("id = ?", achievementId).First(&achievements).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Achievement not found"})
		return
	}

	// tentuin file path dari file yg mau didelete
	filePath := filepath.Join("uploads", achievements.Foto)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, still delete the achievements from database
		if err := ac.DB.Delete(&achievements).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievements from database"})
			return
		}
		ctx.JSON(http.StatusOK, "Achievement deleted successfully")
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete achievement from database
	if err := ac.DB.Delete(&achievements).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete achievements from database"})
		return
	}

	ctx.JSON(http.StatusOK, "Achievement deleted successfully")
}
