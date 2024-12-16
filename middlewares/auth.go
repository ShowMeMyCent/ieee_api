package middlewares

import (
	"backend/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware - middleware untuk verifikasi JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Ambil header Authorization
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			// Jika tidak ada header Authorization, kembalikan error
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		// Ekstrak token dari header Authorization
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == "" {
			// Jika token tidak ada setelah prefix "Bearer ", kembalikan error
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is missing"})
			ctx.Abort()
			return
		}

		// Validasi token dan ambil userID dari token
		userID, err := utils.ValidateToken(tokenString)
		if err != nil {
			// Jika token tidak valid atau ada masalah, kembalikan error
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token", "details": err.Error()})
			ctx.Abort()
			return
		}

		// Attach userID ke context, supaya bisa diakses di handler selanjutnya
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
