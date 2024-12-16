package utils

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
)

// Secret key untuk menandatangani token
var secretKey = []byte("your_secret_key")

// ValidateToken akan memverifikasi dan mengembalikan userID dari JWT yang valid
func ValidateToken(tokenString string) (uint, error) {
	// Parse token menggunakan secret key
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validasi bahwa token menggunakan algoritma HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return secretKey, nil
	})

	// Jika terjadi error parsing, kembalikan error
	if err != nil {
		return 0, err
	}

	// Pastikan token valid
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Mengambil userID dari claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("user_id not found in token")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}
