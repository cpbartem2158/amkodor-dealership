package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type EmployeeService struct {
	repo repository.EmployeeRepository
}

func NewEmployeeService(repo repository.EmployeeRepository) *EmployeeService {
	return &EmployeeService{repo: repo}
}

func (s *EmployeeService) GetAll(limit, offset int) ([]models.Employee, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *EmployeeService) GetByID(id int) (*models.Employee, error) {
	return s.repo.GetByID(id)
}

func (s *EmployeeService) Create(employee *models.Employee) (int, error) {
	return s.repo.Create(employee)
}

func (s *EmployeeService) Update(employee *models.Employee) error {
	return s.repo.Update(employee)
}

func (s *EmployeeService) Delete(id int) error {
	return s.repo.Delete(id)
}