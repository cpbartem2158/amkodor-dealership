package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type ServiceOrderService struct {
	repo *repository.ServiceRepository
}

func NewServiceOrderService(repo *repository.ServiceRepository) *ServiceOrderService {
	return &ServiceOrderService{repo: repo}
}

func (s *ServiceOrderService) GetAllOrders(limit, offset int) ([]models.ServiceOrder, error) {
	return s.repo.GetAllOrders(limit, offset)
}

func (s *ServiceOrderService) GetOrderByID(id int) (*models.ServiceOrder, error) {
	return s.repo.GetOrderByID(id)
}

func (s *ServiceOrderService) CreateOrder(so *models.ServiceOrder) (int, error) {
	return s.repo.CreateOrder(so)
}

func (s *ServiceOrderService) UpdateOrder(so *models.ServiceOrder) error {
	return s.repo.UpdateOrder(so)
}

func (s *ServiceOrderService) CompleteOrder(orderID int) error {
	return s.repo.CompleteOrder(orderID)
}

func (s *ServiceOrderService) GetAllParts(limit, offset int) ([]models.SparePart, error) {
	return s.repo.GetAllParts(limit, offset)
}

func (s *ServiceOrderService) GetPartByID(id int) (*models.SparePart, error) {
	return s.repo.GetPartByID(id)
}

func (s *ServiceOrderService) CreatePart(sp *models.SparePart) (int, error) {
	return s.repo.CreatePart(sp)
}

func (s *ServiceOrderService) UpdatePart(sp *models.SparePart) error {
	return s.repo.UpdatePart(sp)
}

func (s *ServiceOrderService) GetAllTestDrives(limit, offset int) ([]models.TestDrive, error) {
	return s.repo.GetAllTestDrives(limit, offset)
}

func (s *ServiceOrderService) CreateTestDrive(td *models.TestDrive) (int, error) {
	return s.repo.CreateTestDrive(td)
}

func (s *ServiceOrderService) UpdateTestDrive(td *models.TestDrive) error {
	return s.repo.UpdateTestDrive(td)
}
