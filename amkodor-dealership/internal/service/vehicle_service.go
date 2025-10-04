package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

// VehicleService - сервис для работы с техникой
type VehicleService struct {
	repo *repository.VehicleRepository
}

func NewVehicleService(repo *repository.VehicleRepository) *VehicleService {
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GetAll(limit, offset int) ([]models.Vehicle, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *VehicleService) GetByID(id int) (*models.Vehicle, error) {
	return s.repo.GetByID(id)
}

func (s *VehicleService) Create(v *models.Vehicle) (int, error) {
	return s.repo.Create(v)
}

func (s *VehicleService) Update(v *models.Vehicle) error {
	return s.repo.Update(v)
}

func (s *VehicleService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *VehicleService) Search(params map[string]interface{}) ([]models.Vehicle, error) {
	return s.repo.Search(params)
}

func (s *VehicleService) GetHistory(vehicleID int) ([]models.VehicleHistory, error) {
	return s.repo.GetHistory(vehicleID)
}
