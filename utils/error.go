package utils

import (
	"github.com/gin-gonic/gin"
)

// HandleError - helper untuk mengirim error response
func HandleError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"error": message,
	})
}
