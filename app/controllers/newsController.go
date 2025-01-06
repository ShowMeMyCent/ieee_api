package controllers

import (
	"backend/app/models"
	"backend/app/services"
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type NewsResponse struct {
	ID          uint        `json:"id"`
	Title       string      `json:"title"`
	Kategori    string      `json:"kategori"`
	Thumbnail   string      `json:"thumbnail"`
	IsiKonten   interface{} `json:"isi_konten"`
	NamaPenulis string      `json:"nama_penulis"`
	Link        string      `json:"link"`
	ImageURL    string      `json:"image_url"`
	Date        string      `json:"date"`
}

type NewsController struct {
	DB      *gorm.DB
	Service *services.NewsService
}

// GetAllNews adalah fungsi untuk mendapatkan semua news dari database.
// @Summary Get All News
// @Description Retrieves all news from the database.
// @Tags News
// @Produce json
// @Success 200 {array} models.News
// @Router /news [get]
func (nc *NewsController) GetAllNews(ctx *gin.Context) {
	news, err := nc.Service.GetAllNews()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get news from database"})
		return
	}

	var response []NewsResponse
	for _, n := range news {
		var isiKontenDecoded interface{}
		if err := json.Unmarshal([]byte(n.IsiKonten), &isiKontenDecoded); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode isi_konten"})
			return
		}

		response = append(response, NewsResponse{
			ID:          n.ID,
			Title:       n.Title,
			Kategori:    n.Kategori,
			Thumbnail:   n.Thumbnail,
			IsiKonten:   isiKontenDecoded,
			NamaPenulis: n.NamaPenulis,
			Link:        n.Link,
			ImageURL:    n.ImageURL,
			Date:        n.Date,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// GetNewsById adalah fungsi untuk mengambil news berdasarkan ID.
// @Summary Get News By ID
// @Description Retrieves news data by its ID.
// @Tags News
// @Param id path string true "News ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /news/{id} [get]
func (nc *NewsController) GetNewsById(ctx *gin.Context) {
	newsId := ctx.Param("id")
	news, err := nc.Service.GetNewsById(newsId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	var isiKontenDecoded interface{}
	if err := json.Unmarshal([]byte(news.IsiKonten), &isiKontenDecoded); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode isi_konten"})
		return
	}

	response := NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Kategori:    news.Kategori,
		Thumbnail:   news.Thumbnail,
		IsiKonten:   isiKontenDecoded,
		NamaPenulis: news.NamaPenulis,
		Link:        news.Link,
		ImageURL:    news.ImageURL,
		Date:        news.Date,
	}

	ctx.JSON(http.StatusOK, response)
}

// GetThumbnailNews adalah fungsi untuk mengambil Thumbnail News berdasarkan ID.
// @Summary Get Thumbnail News
// @Description Retrieves the image of a News by its ID.
// @Tags News
// @Param id path string true "News ID"
// @Produce octet-stream
// @Success 200 {file} octet-stream
// @Router /news/file/{id} [get]
func (nc *NewsController) GetThumbnailNews(ctx *gin.Context) {
	newsID := ctx.Param("id")
	news, err := nc.Service.GetNewsById(newsID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	filePath := filepath.Join("uploads", news.Thumbnail)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	ctx.Header("Content-Description", "File Transfer")
	ctx.Header("Content-Transfer-Encoding", "binary")
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", news.Thumbnail))

	ext := filepath.Ext(news.Thumbnail)
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

// GetNewsByCategory adalah fungsi untuk mengambil news berdasarkan Category.
// @Summary Get All News By Category
// @Description Retrieves news data by its Category.
// @Tags News
// @Param category path string true "News Category"
// @Produce json
// @Success 200 {array} NewsResponse
// @Router /news/category/{category} [get]
func (nc *NewsController) GetNewsByCategory(ctx *gin.Context) {
	newsCategory := ctx.Param("category")
	newsList, err := nc.Service.GetNewsByCategory(newsCategory)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	var response []NewsResponse
	for _, news := range newsList {
		var isiKontenDecoded interface{}
		if err := json.Unmarshal([]byte(news.IsiKonten), &isiKontenDecoded); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode isi_konten"})
			return
		}

		response = append(response, NewsResponse{
			ID:          news.ID,
			Title:       news.Title,
			Kategori:    news.Kategori,
			Thumbnail:   news.Thumbnail,
			IsiKonten:   isiKontenDecoded,
			NamaPenulis: news.NamaPenulis,
			Link:        news.Link,
			ImageURL:    news.ImageURL,
			Date:        news.Date,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

// InsertNews adalah fungsi untuk membuat post news terbaru.
// @Summary Insert a new news
// @Description Insert a news and saves them to the database.
// @Tags News
// @Accept multipart/form-data
// @Param title formData string true "Judul news"
// @Param kategori formData string true "Kategori news"
// @Param thumbnail formData file true "Thumbnail news"
// @Param isi_konten formData string true "Isi konten news"
// @Param nama_penulis formData string true "Nama penulis news"
// @Param link formData string true "Link news"
// @Param date formData string true "Date news"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.News
// @Router /news [post]
func (nc *NewsController) InsertNews(ctx *gin.Context) {
	file, err := ctx.FormFile("thumbnail")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	host := utils.Getenv("ENV_HOST", "localhost")

	fileName := uuid.New().String() + filepath.Ext(file.Filename)
	filePath := filepath.Join("uploads", fileName)

	var news models.News
	news.Title = ctx.PostForm("title")
	news.Kategori = ctx.PostForm("kategori")
	news.Date = ctx.PostForm("date")
	isiKonten := ctx.PostForm("isi_konten")

	var isiKontenJSON interface{}
	if err := json.Unmarshal([]byte(isiKonten), &isiKontenJSON); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format for isi_konten"})
		return
	}

	isiKontenBytes, err := json.Marshal(isiKontenJSON)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process isi_konten"})
		return
	}
	news.IsiKonten = string(isiKontenBytes)
	news.NamaPenulis = ctx.PostForm("nama_penulis")
	news.Link = ctx.PostForm("link")
	news.Thumbnail = fileName
	news.ImageURL = fmt.Sprintf("https://%s/uploads/%s", host, fileName)

	validationErrors := utils.ValidateStruct(news)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	if err := nc.Service.CreateNews(&news); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save news to database"})
		return
	}

	var isiKontenDecoded interface{}
	if err := json.Unmarshal([]byte(news.IsiKonten), &isiKontenDecoded); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode isi_konten"})
		return
	}

	response := NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Kategori:    news.Kategori,
		Thumbnail:   news.Thumbnail,
		IsiKonten:   isiKontenDecoded,
		NamaPenulis: news.NamaPenulis,
		Link:        news.Link,
		ImageURL:    news.ImageURL,
		Date:        news.Date,
	}

	ctx.JSON(http.StatusOK, response)
}

// EditNews adalah fungsi untuk mengedit News
// @Summary Edit News
// @Description Edits a News by its ID
// @Tags News
// @Accept multipart/form-data
// @Param id path string true "News ID"
// @Param title formData string true "Title News"
// @Param kategori formData string true "Kategori News"
// @Param thumbnail formData file false "Thumbnail News"
// @Param isi_konten formData string true "Isi Konten News"
// @Param nama_penulis formData string true "Nama Penulis News"
// @Param link formData string true "Link News"
// @Param date formData string true "Date News"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {object} models.News
// @Router /news/{id} [put]

func (nc *NewsController) EditNews(ctx *gin.Context) {
	newsID := ctx.Param("id")
	host := utils.Getenv("ENV_HOST", "localhost")

	news, err := nc.Service.GetNewsById(newsID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	news.Title = ctx.PostForm("title")
	news.Kategori = ctx.PostForm("kategori")
	news.Date = ctx.PostForm("date")
	isiKonten := ctx.PostForm("isi_konten")

	if isiKonten != "" {
		var isiKontenJSON interface{}
		if err := json.Unmarshal([]byte(isiKonten), &isiKontenJSON); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format for isi_konten"})
			return
		}

		isiKontenBytes, err := json.Marshal(isiKontenJSON)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process isi_konten"})
			return
		}
		news.IsiKonten = string(isiKontenBytes)
	}
	news.NamaPenulis = ctx.PostForm("nama_penulis")
	news.Link = ctx.PostForm("link")

	validationErrors := utils.ValidateStruct(news)
	if len(validationErrors) > 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErrors})
		return
	}

	file, err := ctx.FormFile("thumbnail")
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

		oldFilePath := filepath.Join("uploads", news.Thumbnail)
		if err := os.Remove(oldFilePath); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove old file"})
			return
		}

		news.Thumbnail = fileName
		news.ImageURL = fmt.Sprintf("https://%s/uploads/%s", host, fileName)
	}

	if err := nc.Service.UpdateNews(news); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save updated news to database"})
		return
	}

	var isiKontenDecoded interface{}
	if err := json.Unmarshal([]byte(news.IsiKonten), &isiKontenDecoded); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode isi_konten"})
		return
	}

	response := NewsResponse{
		ID:          news.ID,
		Title:       news.Title,
		Kategori:    news.Kategori,
		Thumbnail:   news.Thumbnail,
		IsiKonten:   isiKontenDecoded,
		NamaPenulis: news.NamaPenulis,
		Link:        news.Link,
		ImageURL:    news.ImageURL,
		Date:        news.Date,
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteNews adalah fungsi untuk menghapus News dan gambar-nya dari database.
// @Summary Delete News
// @Description Delete a News and its thumbnail from the database and storage.
// @Tags News
// @Param id path string true "News ID"
// @Param Authorization header string true "Authorization. How to input in swagger : 'Bearer <insert_your_token_here>'"
// @Produce json
// @Success 200 {string} string "News deleted successfully"
// @Router /news/{id} [delete]
func (nc *NewsController) DeleteNews(ctx *gin.Context) {
	newsId := ctx.Param("id")

	news, err := nc.Service.GetNewsById(newsId)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "News not found"})
		return
	}

	filePath := filepath.Join("uploads", news.Thumbnail)
	_, err = os.Stat(filePath)
	if os.IsNotExist(err) {
		if err := nc.Service.DeleteNewsById(newsId); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news from database"})
			return
		}
		ctx.JSON(http.StatusOK, "News deleted successfully")
		return
	}

	if err := os.Remove(filePath); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file"})
		return
	}

	if err := nc.Service.DeleteNewsById(newsId); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete news from database"})
		return
	}

	ctx.JSON(http.StatusOK, "News deleted successfully")
}
