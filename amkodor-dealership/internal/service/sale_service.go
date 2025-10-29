package service

import (
	"context"
	"database/sql"
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"fmt"
)

type SaleService struct {
	repo        repository.SaleRepository
	vehicleRepo repository.VehicleRepository
}

func NewSaleService(repo repository.SaleRepository, vehicleRepo repository.VehicleRepository) *SaleService {
	return &SaleService{repo: repo, vehicleRepo: vehicleRepo}
}

func (s *SaleService) GetAll(limit, offset int) ([]models.Sale, error) {
	return s.repo.GetAll(context.Background())
}

func (s *SaleService) GetByID(id int) (*models.Sale, error) {
	return s.repo.GetByID(context.Background(), id)
}

func (s *SaleService) Create(vehicleID int, customerID, corporateClientID *int, employeeID int,
	paymentType string, additionalDiscount float64, contractNumber, notes string) (int, error) {

	// Проверка доступности техники
	vehicle, err := s.vehicleRepo.GetByID(context.Background(), vehicleID)
	if err != nil {
		return 0, fmt.Errorf("vehicle not found: %w", err)
	}

	if vehicle.Status != "В наличии" && vehicle.Status != "Зарезервировано" {
		return 0, fmt.Errorf("vehicle is not available for sale")
	}

	sale := &models.Sale{
		VehicleID:  vehicleID,
		EmployeeID: employeeID,
		PaymentType: paymentType,
	}
	
	if customerID != nil {
		sale.CustomerID = sql.NullInt64{Int64: int64(*customerID), Valid: true}
	}
	
	if corporateClientID != nil {
		sale.CorporateClientID = sql.NullInt64{Int64: int64(*corporateClientID), Valid: true}
	}
	
	if contractNumber != "" {
		sale.ContractNumber = sql.NullString{String: contractNumber, Valid: true}
	}
	
	if notes != "" {
		sale.Notes = sql.NullString{String: notes, Valid: true}
	}
	
	return s.repo.Create(context.Background(), sale)
}

func (s *SaleService) Update(sale *models.Sale) error {
	return s.repo.Update(context.Background(), sale)
}

func (s *SaleService) Delete(id int) error {
	return s.repo.Cancel(context.Background(), id)
}

func (s *SaleService) GetHistory(saleID int) ([]models.SaleHistory, error) {
	return s.repo.GetHistory(context.Background(), saleID)
}
