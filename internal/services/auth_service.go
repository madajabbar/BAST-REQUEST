package services

import (
	"bast-request/internal/models"
	"bast-request/internal/repositories"
	"bast-request/internal/utils"
	"fmt"
)

type AuthService struct {
	authRepo *repositories.AuthRepository
}

func NewAuthService(authRepo *repositories.AuthRepository) *AuthService {
	return &AuthService{authRepo: authRepo}
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.authRepo.FindUserByEmail(email)
	if err != nil {
		return "", fmt.Errorf("invalid email or password") // do not expose whether user exists
	}
	
	if !utils.ValidatePassword(user.PasswordHash, password) {
		return "", fmt.Errorf("invalid email or password")
	}
	
	// Assuming Role is preloaded or we just use RoleID. Wait, GenerateToken needs role name.
	// We need to fetch role name if it's not preloaded. But in this case, we might just pass the role name string.
	// Actually utils.GenerateToken expects role name as string.
	// If user.Role.Name is empty because it wasn't preloaded, we might need to fetch it.
	// Or we can modify FindUserByEmail to preload Role. Let's assume we can fetch role if not preloaded.
	// Actually in models.User, Role is a struct. Let's just pass user.Role.Name for now, but we should make sure FindUserByEmail preloads it.
	// I'll update auth_repository.go to preload Role.

	token, err := utils.GenerateToken(user.UserID.String(), user.Role.Name)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *AuthService) Register(req *models.RegisterRequest) (string, error) {
	// Check if user already exists
	_, err := s.authRepo.FindUserByEmail(req.Email)
	if err == nil {
		return "", fmt.Errorf("user already exists")
	}

	// Fetch Role ID based on Role Name
	role, err := s.authRepo.FindRoleByName(req.Role)
	if err != nil {
		return "", fmt.Errorf("invalid role: %v", err)
	}

	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return "", err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		RoleID:       role.RoleID,
	}

	err = s.authRepo.CreateUser(user)
	if err != nil {
		return "", err
	}
	return "User registered successfully", nil
}
