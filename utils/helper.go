package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ga usah dipake gpp si cuma enak aja klo ada
func Getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return ""
	}
}

func ValidateStruct(data interface{}) map[string]string {
	// Create a new validator instance
	validate := validator.New()

	// Map to store validation errors
	validationErrors := make(map[string]string)

	// Validate the data
	if err := validate.Struct(data); err != nil {
		// Type assert to get the validation errors
		for _, err := range err.(validator.ValidationErrors) {
			// Extract field name from struct tag
			fieldName := strings.ToLower(err.Field())

			// Store validation error
			validationErrors[fieldName] = err.Tag()
		}
	}

	return validationErrors
}
