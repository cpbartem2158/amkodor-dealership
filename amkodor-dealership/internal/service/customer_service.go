package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type CustomerService struct {
	repo repository.CustomerRepository
}

func NewCustomerService(repo repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetAll(limit, offset int) ([]models.Customer, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *CustomerService) GetByID(id int) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) Create(customer *models.Customer) (int, error) {
	return s.repo.Create(customer)
}

func (s *CustomerService) Update(customer *models.Customer) error {
	return s.repo.Update(customer)
}

func (s *CustomerService) Delete(id int) error {
	return s.repo.Delete(id)
}