package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/utils"
	"fmt"
)

type EmployeeService struct {
	repo *repository.EmployeeRepository
}

func NewEmployeeService(repo *repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAll(limit, offset int) ([]models.Employee, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *EmployeeService) GetByID(id int) (*models.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) Create(e *models.Employee, password string) (int, error) {
	// Хеширование пароля
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return 0, fmt.Errorf("error hashing password: %w", err)
	}

	e.PasswordHash = passwordHash
	return s.repo.Create(e)
}

func (s *EmployeeService) Update(e *models.Employee) error {
	return s.repo.Update(e)
}

func (s *EmployeeService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *EmployeeService) UpdatePassword(employeeID int, newPassword string) error {
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	return s.repo.UpdatePassword(employeeID, passwordHash)
}
