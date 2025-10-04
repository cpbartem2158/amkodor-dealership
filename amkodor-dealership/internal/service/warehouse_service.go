package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type WarehouseService struct {
	repo *repository.WarehouseRepository
}

func NewWarehouseService(repo *repository.WarehouseRepository) *WarehouseService {
	return &WarehouseService{repo: repo}
}

func (s *WarehouseService) GetAll() ([]models.Warehouse, error) {
	return s.repo.GetAll()
}

func (s *WarehouseService) GetByID(id int) (*models.Warehouse, error) {
	return s.repo.GetByID(id)
}

func (s *WarehouseService) Create(w *models.Warehouse) (int, error) {
	return s.repo.Create(w)
}

func (s *WarehouseService) Update(w *models.Warehouse) error {
	return s.repo.Update(w)
}

func (s *WarehouseService) GetStatistics(id int) (*models.Warehouse, error) {
	return s.repo.GetStatistics(id)
}
