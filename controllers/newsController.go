package controllers

import (
	"backend/models"
	"backend/services"
	"backend/utils"
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type NewsController struct {
	NewsService services.NewsService
}

func NewNewsController(newsService services.NewsService) *NewsController {
	return &NewsController{NewsService: newsService}
}

// GetAllNews mengambil seluruh berita
func (nc *NewsController) GetAllNews(ctx *gin.Context) {
	// Set timeout untuk konteks request
	c, cancel := context.WithTimeout(ctx.Request.Context(), 10*time.Second)
	defer cancel()

	news, err := nc.NewsService.GetAllNews(c)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, news)
}

// InsertNews menyimpan berita baru beserta thumbnail
func (nc *NewsController) InsertNews(ctx *gin.Context) {
	// Set timeout untuk konteks request
	c, cancel := context.WithTimeout(ctx.Request.Context(), 15*time.Second)
	defer cancel()

	// Validasi file input
	file, err := ctx.FormFile("thumbnail")
	if err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Thumbnail is required")
		return
	}

	// Validasi data form lainnya
	news := models.News{
		Title:       ctx.PostForm("title"),
		Kategori:    ctx.PostForm("kategori"),
		Date:        ctx.PostForm("date"),
		IsiKonten:   ctx.PostForm("isi_konten"),
		NamaPenulis: ctx.PostForm("nama_penulis"),
		Link:        ctx.PostForm("link"),
	}

	// Validasi struct berita
	validationErrors := utils.ValidateStruct(news)
	if len(validationErrors) > 0 {
		// Menggabungkan kesalahan validasi menjadi satu string
		var errorMessages []string
		for _, message := range validationErrors {
			errorMessages = append(errorMessages, message)
		}
		utils.HandleError(ctx, http.StatusBadRequest, strings.Join(errorMessages, ", "))
		return
	}

	// Panggil service untuk menyimpan berita
	response, err := nc.NewsService.InsertNews(c, news, file, ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "News successfully inserted",
		"data":    response,
	})
}
