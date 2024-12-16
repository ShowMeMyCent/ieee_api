package middlewares

import (
	"backend/utils/token"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           7 * time.Hour,
	})
}

func AdminCheckMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := token.TokenValid(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Setelah token divalidasi, kita dapat memeriksa role dari user yang terautentikasi.
		role, err := token.ExtractUserRole(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		// Jika role tidak sesuai, berikan pesan error dan hentikan proses.
		if role != "admin" {
			c.String(http.StatusForbidden, "Access denied. Insufficient role.")
			c.Abort()
			return
		}
	}
}
