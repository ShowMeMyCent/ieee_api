package services

import (
	"backend/repositories"
	"errors"
)

// AuthAdminService - interface untuk autentikasi admin
type AuthAdminService interface {
	LoginAdmin(kode string) (string, error)
}

type authAdminService struct {
	repo repositories.UserRepository
}

// NewAuthAdminService - inisialisasi AuthAdminService baru
func NewAuthAdminService(repo repositories.UserRepository) AuthAdminService {
	return &authAdminService{repo} // Mengembalikan instance authAdminService sebagai AuthAdminService
}

// LoginAdmin - implementasi autentikasi admin
func (s *authAdminService) LoginAdmin(kode string) (string, error) {
	// Cari user berdasarkan kode
	user, err := s.repo.FindByCode(kode)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate token JWT
	token, err := user.GenerateJWT()
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}
