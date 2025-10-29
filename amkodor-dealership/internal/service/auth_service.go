package service

import (
	"context"
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/utils"
	"fmt"
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

func (s *AuthService) Login(email, password string) (*models.User, string, error) {
	// Получение пользователя по email
	user, err := s.userRepo.GetByEmail(context.Background(), email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Проверка пароля
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	// Генерация JWT токена
	token, err := utils.GenerateJWT(user.UserID, email, s.jwtSecret, 24)
	if err != nil {
		return nil, "", fmt.Errorf("error generating token: %w", err)
	}

	return user, token, nil
}
