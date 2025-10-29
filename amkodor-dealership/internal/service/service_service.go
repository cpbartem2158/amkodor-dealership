package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type ServiceService struct {
	repo repository.ServiceRepository
}

func NewServiceService(repo repository.ServiceRepository) *ServiceService {
	return &ServiceService{repo: repo}
}

func (s *ServiceService) GetAll(limit, offset int) ([]models.ServiceOrder, error) {
	return s.repo.GetAllOrders(limit, offset)
}

func (s *ServiceService) GetByID(id int) (*models.ServiceOrder, error) {
	return s.repo.GetByID(id)
}

func (s *ServiceService) Create(order *models.ServiceOrder) (int, error) {
	return s.repo.Create(order)
}

func (s *ServiceService) Update(order *models.ServiceOrder) error {
	return s.repo.Update(order)
}

func (s *ServiceService) Delete(id int) error {
	return s.repo.Delete(id)
}
