package service

import (
	"fmt"
	"strings"

	"github.com/emirrasyad/loopaffi-backend/internal/dto"
	"github.com/emirrasyad/loopaffi-backend/internal/entity"
	"github.com/emirrasyad/loopaffi-backend/internal/repository"
	"github.com/emirrasyad/loopaffi-backend/internal/utils"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal mencari user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("email atau password salah")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, fmt.Errorf("email atau password salah")
	}

	// Gunakan RoleID untuk token
	token, err := utils.GenerateToken(user.ID, user.Email, user.RoleID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserProfile{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID, // Kembalikan RoleID (misal: role_admin atau role_affiliate)
		},
	}, nil
}

func (s *AuthService) Register(req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Check if email already exists
	existing, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("gagal memeriksa email: %w", err)
	}
	if existing != nil {
		return nil, fmt.Errorf("email sudah terdaftar")
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	// Tentukan role berdasarkan email (Hack untuk dev: jika ada kata admin, jadi admin)
	roleID := "affiliate"
	if strings.Contains(strings.ToLower(req.Email), "admin") {
		roleID = "admin"
	}

	user := &entity.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		RoleID:       roleID,
		Phone:        req.Whatsapp,
		Status:       "active",
	}

	// Create user (sudah include info affiliate di tabel yang sama)
	err = s.userRepo.CreateUser(user)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat akun: %w", err)
	}

	// Generate token so user is auto-logged-in after register
	token, err := utils.GenerateToken(user.ID, user.Email, user.RoleID, s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat token: %w", err)
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserProfile{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.RoleID,
		},
	}, nil
}

func (s *AuthService) ForgotPassword(req dto.ForgotPasswordRequest) error {
	email := strings.TrimSpace(strings.ToLower(req.Email))

	if len(req.NewPassword) < 6 {
		return fmt.Errorf("password minimal 6 karakter")
	}

	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return fmt.Errorf("gagal mencari user: %w", err)
	}

	if user == nil {
		return fmt.Errorf("email tidak ditemukan")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	err = s.userRepo.UpdatePasswordByEmail(email, hashedPassword)
	if err != nil {
		return fmt.Errorf("gagal mengubah password: %w", err)
	}

	return nil
}
