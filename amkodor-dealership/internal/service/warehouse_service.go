package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type WarehouseService struct {
	repo repository.WarehouseRepository
}

func NewWarehouseService(repo repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) GetAll() ([]models.Warehouse, error) {
	return s.repo.GetAll()
}

func (s *WarehouseService) GetByID(id int) (*models.Warehouse, error) {
	return s.repo.GetByID(id)
}

func (s *WarehouseService) Create(warehouse *models.Warehouse) (int, error) {
	return s.repo.Create(warehouse)
}

func (s *WarehouseService) Update(warehouse *models.Warehouse) error {
	return s.repo.Update(warehouse)
}

func (s *WarehouseService) Delete(id int) error {
	return s.repo.Delete(id)
}