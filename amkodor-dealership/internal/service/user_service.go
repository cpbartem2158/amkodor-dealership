package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/utils"
	"context"
	"fmt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Register регистрирует нового пользователя
func (s *UserService) Register(ctx context.Context, name, email, phone, password string) (*models.User, error) {
	// Проверяем, что email не занят
	exists, err := s.repo.EmailExists(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("error checking email: %w", err)
	}

	if exists {
		return nil, fmt.Errorf("user with email %s already exists", email)
	}

	// Хешируем пароль
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("error hashing password: %w", err)
	}

	// Создаем пользователя
	user := &models.User{
		Name:         name,
		Email:        email,
		Phone:        phone,
		PasswordHash: hashedPassword,
		Role:         "user", // По умолчанию обычный пользователь
	}

	userID, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}

	user.UserID = userID
	// Не возвращаем хеш пароля
	user.PasswordHash = ""

	return user, nil
}

// Login авторизует пользователя
func (s *UserService) Login(ctx context.Context, email, password string) (*models.User, error) {
	// Получаем пользователя по email
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Проверяем пароль
	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""

	return user, nil
}

// GetByID получает пользователя по ID
func (s *UserService) GetByID(ctx context.Context, userID int) (*models.User, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""

	return user, nil
}

// GetUserByToken получает пользователя по JWT токену
func (s *UserService) GetUserByToken(ctx context.Context, tokenString string) (*models.User, error) {
	// Валидируем токен
	claims, err := utils.ValidateJWTClaims(tokenString, "amkodor-secret-key-change-in-production")
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Получаем пользователя по ID из токена
	user, err := s.repo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Не возвращаем хеш пароля
	user.PasswordHash = ""

	return user, nil
}

// UpdateProfile обновляет профиль пользователя
func (s *UserService) UpdateProfile(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

// Update обновляет данные пользователя
func (s *UserService) Update(ctx context.Context, user *models.User) error {
	return s.repo.Update(ctx, user)
}

// Delete удаляет пользователя
func (s *UserService) Delete(ctx context.Context, userID int) error {
	return s.repo.Delete(ctx, userID)
}
