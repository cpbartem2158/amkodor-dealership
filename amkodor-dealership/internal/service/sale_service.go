package service

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"fmt"
)

type SaleService struct {
	repo        *repository.SaleRepository
	vehicleRepo *repository.VehicleRepository
}

func NewSaleService(repo *repository.SaleRepository, vehicleRepo *repository.VehicleRepository) *SaleService {
	return &SaleService{repo: repo, vehicleRepo: vehicleRepo}
}

func (s *SaleService) GetAll(limit, offset int) ([]models.Sale, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *SaleService) GetByID(id int) (*models.Sale, error) {
	return s.repo.GetByID(id)
}

func (s *SaleService) Create(vehicleID int, customerID, corporateClientID *int, employeeID int,
	paymentType string, additionalDiscount float64, contractNumber, notes string) (int, error) {

	// Проверка доступности техники
	vehicle, err := s.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		return 0, fmt.Errorf("vehicle not found: %w", err)
	}

	if vehicle.Status != "В наличии" && vehicle.Status != "Зарезервировано" {
		return 0, fmt.Errorf("vehicle is not available for sale")
	}

	return s.repo.Create(vehicleID, customerID, corporateClientID, employeeID,
		paymentType, additionalDiscount, contractNumber, notes)
}

func (s *SaleService) Update(sale *models.Sale) error {
	return s.repo.Update(sale)
}

func (s *SaleService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *SaleService) GetHistory(saleID int) ([]models.SaleHistory, error) {
	return s.repo.GetHistory(saleID)
}
