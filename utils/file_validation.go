package utils

import (
	"fmt"
	"mime/multipart"
)

const maxSize = 5 * 1024 * 1024 // 5MB

// ValidateFile checks if the uploaded file is valid (image type, size)
func ValidateFile(file *multipart.FileHeader) error {
	// Check file size
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds the limit of 5MB")
	}

	// Check file type
	allowedTypes := []string{"image/png", "image/jpeg", "image/jpg"}
	if !Contains(allowedTypes, file.Header.Get("Content-Type")) {
		return fmt.Errorf("invalid file type. Only PNG, JPG, and JPEG are allowed")
	}

	return nil
}

// Contains checks if a value is present in a slice
func Contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}
