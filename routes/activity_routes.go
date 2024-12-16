package routes

import (
	"backend/controllers"
	"github.com/gin-gonic/gin"
)

// RegisterActivityRoutes untuk mendaftarkan semua routes terkait aktivitas
func RegisterActivityRoutes(router *gin.Engine, activitiesController *controllers.ActivitiesController) {
	activity := router.Group("/activities") // grup route untuk activities
	{
		// Route untuk upload aktivitas (POST)
		activity.POST("/upload", activitiesController.UploadActivity)

		// Route untuk mendapatkan semua aktivitas (GET)
		activity.GET("/", activitiesController.GetAllActivities)

		// Route untuk mendapatkan aktivitas berdasarkan ID (GET)
		activity.GET("/:id", activitiesController.GetActivityByID)

		// Route untuk mengedit aktivitas (PUT)
		activity.PUT("/:id", activitiesController.EditActivity)

		// Route untuk menghapus aktivitas (DELETE)
		activity.DELETE("/:id", activitiesController.DeleteActivity)
	}
}
