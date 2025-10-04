package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) GetAll(limit, offset int) ([]models.Customer, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *CustomerService) GetByID(id int) (*models.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) Create(c *models.Customer) (int, error) {
	return s.repo.Create(c)
}

func (s *CustomerService) Update(c *models.Customer) error {
	return s.repo.Update(c)
}

func (s *CustomerService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *CustomerService) Search(params map[string]interface{}) ([]models.Customer, error) {
	return s.repo.Search(params)
}

func (s *CustomerService) GetAllCorporate(limit, offset int) ([]models.CorporateClient, error) {
	return s.repo.GetAllCorporate(limit, offset)
}

func (s *CustomerService) GetCorporateByID(id int) (*models.CorporateClient, error) {
	return s.repo.GetCorporateByID(id)
}

func (s *CustomerService) CreateCorporate(c *models.CorporateClient) (int, error) {
	return s.repo.CreateCorporate(c)
}

func (s *CustomerService) UpdateCorporate(c *models.CorporateClient) error {
	return s.repo.UpdateCorporate(c)
}

func (s *CustomerService) DeleteCorporate(id int) error {
	return s.repo.DeleteCorporate(id)
}
