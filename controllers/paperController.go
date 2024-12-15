package controllers

import (
	"backend/models"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"path/filepath"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaperController struct {
	DB *gorm.DB
}

// UploadPaper adalah fungsi untuk mengupload paper beserta file-nya.
// @Summary Upload Paper with File
// @Description Uploads a paper along with its file and saves them to the database.
// @Tags Papers
// @Accept multipart/form-data
// @Param judul formData string true "Judul paper"
// @Param abstrak formData string true "Abstrak paper"
// @Param link formData string true "Link paper"
// @Param file_paper formData file true "File paper"
// @Param author formData string true "Author paper"
// @Param tanggal_terbit formData string true "Tanggal terbit paper"
// @Produce json
// @Success 200 {object} models.Paper
// @Router /papers [post]
func (pc *PaperController) UploadPaper(ctx *gin.Context) {
	// Get the file from the form data
	file, err := ctx.FormFile("file_paper")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Define the path where the file will be saved, pake UUID, untuk skg taro di uploads dlu
	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("uploads", fileName)

	// Save the file to the defined path
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	judul := ctx.PostForm("judul")
	abstrak := ctx.PostForm("abstrak")
	link := ctx.PostForm("link")
	author := ctx.PostForm("author")
	tanggalTerbit := ctx.PostForm("tanggal_terbit")

	// Create Paper object
	paper := models.Paper{
		Judul:         judul,
		Abstrak:       abstrak,
		Link:          link,
		FilePaper:     fileName,
		Author:        author,
		TanggalTerbit: tanggalTerbit,
	}

	// Save paper to database
	if err := pc.DB.Create(&paper).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save paper to database"})
		return
	}

	ctx.JSON(http.StatusOK, paper)
}

// GetAllPapers adalah fungsi untuk mendapatkan semua paper dari database.
// @Summary Get All Papers
// @Description Retrieves all papers from the database.
// @Tags Papers
// @Produce json
// @Success 200 {array} models.Paper
// @Router /papers [get]
func (pc *PaperController) GetAllPapers(ctx *gin.Context) {
	var papers []models.Paper

	// Retrieve all papers from the database, masalahnyaaa gambarnya ga ngikut gmn ya infoo helppp
	if err := pc.DB.Find(&papers).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get papers from database"})
		return
	}

	// Response with papers
	ctx.JSON(http.StatusOK, papers)
}

// GetPaperFile adalah fungsi untuk mengambil file paper berdasarkan ID.
// @Summary Get Paper File
// @Description Retrieves the file of a paper by its ID.
// @Tags Papers
// @Param id path string true "Paper ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /papers/file/{id} [get]
func (pc *PaperController) GetPaperFile(ctx *gin.Context) {
	// Get paper ID from URL path parameter, pakenya param klo mau diubah ke yg laen sok
	paperID := ctx.Param("id")

	// Retrieve paper from the database by its ID
	var paper models.Paper
	if err := pc.DB.Where("id = ?", paperID).First(&paper).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	// Define the file path
	filePath := filepath.Join("uploads", paper.FilePaper)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Set the headers for the file transfer and return the file
	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", paper.FilePaper))

	// switch case buat content type, sumpah klo ga gini gw gtw gimana biar semua filetype bisa pls help
	ext := filepath.Ext(paper.FilePaper)
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

// EditPaper adalah fungsi untuk mengedit paper, termasuk kemampuan untuk mengganti file paper.
// @Summary Edit Paper
// @Description Edits a paper including the ability to replace its file.
// @Tags Papers
// @Accept multipart/form-data
// @Param id path string true "Paper ID"
// @Param judul formData string true "Judul paper"
// @Param abstrak formData string true "Abstrak paper"
// @Param link formData string true "Link paper"
// @Param file_paper formData file false "File paper (optional)"
// @Param author formData string true "Author paper"
// @Param tanggal_terbit formData string true "Tanggal terbit paper"
// @Produce json
// @Success 200 {object} models.Paper
// @Router /papers/{id} [put]
func (pc *PaperController) EditPaper(ctx *gin.Context) {
	// Get paper ID from URL path parameter
	paperID := ctx.Param("id")

	// Retrieve paper from the database by its ID
	var paper models.Paper
	if err := pc.DB.Where("id = ?", paperID).First(&paper).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	judul := ctx.PostForm("judul")
	abstrak := ctx.PostForm("abstrak")
	link := ctx.PostForm("link")
	author := ctx.PostForm("author")
	tanggalTerbit := ctx.PostForm("tanggal_terbit")

	paper.Judul = judul
	paper.Abstrak = abstrak
	paper.Link = link
	paper.Author = author
	paper.TanggalTerbit = tanggalTerbit

	// Cek apakah file diganti
	file, err := ctx.FormFile("file_paper")
	if err != nil && err != http.ErrMissingFile {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Jika iya, maka di save
	if file != nil {
		// Tentuin tempat nge save
		fileName := uuid.New().String() + filepath.Ext(file.Filename)
		filePath := filepath.Join("uploads", fileName)

		// Save the file to the path
		if err := ctx.SaveUploadedFile(file, filePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			return
		}

		// Remove old file
		oldFilePath := filepath.Join("uploads", paper.FilePaper)
		if err := os.Remove(oldFilePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old file"})
		}

		// Update file paper field in the database
		paper.FilePaper = fileName
	}

	// Save updated paper to database
	if err := pc.DB.Save(&paper).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated paper to database"})
		return
	}

	// Response success
	ctx.JSON(http.StatusOK, paper)
}

// DeletePaper adalah fungsi untuk menghapus paper dan file-nya dari database dan sistem penyimpanan.
// @Summary Delete Paper
// @Description Deletes a paper and its file from the database and storage.
// @Tags Papers
// @Param id path string true "Paper ID"
// @Produce json
// @Success 200 {string} string "Paper deleted successfully"
// @Router /papers/{id} [delete]
func (pc *PaperController) DeletePaper(ctx *gin.Context) {
	paperID := ctx.Param("id")

	// Cari paper dari id
	var paper models.Paper
	if err := pc.DB.Where("id = ?", paperID).First(&paper).Error; err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Paper not found"})
		return
	}

	// tentuin file path dari file yg mau didelete
	filePath := filepath.Join("uploads", paper.FilePaper)

	// Check if the file exists
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// File doesn't exist, still delete the paper from database
		if err := pc.DB.Delete(&paper).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete paper from database"})
			return
		}
		ctx.JSON(http.StatusOK, "Paper deleted successfully")
		return
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	// Delete paper from database
	if err := pc.DB.Delete(&paper).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete paper from database"})
		return
	}

	ctx.JSON(http.StatusOK, "Paper and file deleted successfully")
}
