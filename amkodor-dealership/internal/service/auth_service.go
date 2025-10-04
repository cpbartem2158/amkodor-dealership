package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/utils"
	"fmt"
)

type AuthService struct {
	employeeRepo *repository.EmployeeRepository
	jwtSecret    string
}

func NewAuthService(employeeRepo *repository.EmployeeRepository, jwtSecret string) *AuthService {
	return &AuthService{
		employeeRepo: employeeRepo,
		jwtSecret:    jwtSecret,
	}
}

func (s *AuthService) Login(email, password string) (string, *models.Employee, error) {
	// Получение сотрудника по email
	employee, err := s.employeeRepo.GetByEmail(email)
	if err != nil {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Проверка активности
	if !employee.IsActive {
		return "", nil, fmt.Errorf("account is inactive")
	}

	// Проверка пароля
	if !utils.CheckPasswordHash(password, employee.PasswordHash) {
		return "", nil, fmt.Errorf("invalid credentials")
	}

	// Генерация JWT токена
	token, err := utils.GenerateJWT(employee.EmployeeID, email, s.jwtSecret, 24)
	if err != nil {
		return "", nil, fmt.Errorf("error generating token: %w", err)
	}

	return token, employee, nil
}
