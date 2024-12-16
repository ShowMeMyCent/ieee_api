package utils

//
//import (
//	"fmt"
//	"github.com/dgrijalva/jwt-go"
//)
//
//// ValidateToken validates the JWT token and returns the user ID
//func ValidateToken(tokenString string) (string, error) {
//	// Example secret key, replace with your own secret
//	secretKey := "your-secret-key"
//
//	// Parse the token
//	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//		return []byte(secretKey), nil
//	})
//	if err != nil || !token.Valid {
//		return "", fmt.Errorf("invalid token")
//	}
//
//	// Extract userID from the token claims
//	claims, ok := token.Claims.(jwt.MapClaims)
//	if !ok || !token.Valid {
//		return "", fmt.Errorf("invalid token claims")
//	}
//
//	// Assuming userID is in the claims
//	userID := claims["userID"].(string)
//	return userID, nil
//}
