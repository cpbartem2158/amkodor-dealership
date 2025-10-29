package service

import (
	"amkodor-dealership/internal/repository"
	"database/sql"
)

type Services struct {
	Vehicle          *VehicleService
	Customer         *CustomerService
	Sale             *SaleService
	Employee         *EmployeeService
	Auth             *AuthService
	Dashboard        *DashboardService
	Report           *ReportService
	Admin            *AdminService
	Warehouse        *WarehouseService
	Service          *ServiceService
	ServiceOrderRepo *repository.ServiceOrderRepository
	Favorite         *FavoriteService
}

func NewServices(db *sql.DB, repos *repository.Repository) *Services {
	return &Services{
		Vehicle:          NewVehicleService(&repos.Vehicle),
		Customer:         NewCustomerService(repos.Customer),
		Sale:             NewSaleService(repos.Sale, repos.Vehicle),
		Employee:         NewEmployeeService(repos.Employee),
		Auth:             NewAuthService(&repos.User, "amkodor-secret-key-change-in-production"),
		Dashboard:        NewDashboardService(repos.Dashboard),
		Report:           NewReportService(db),
		Admin:            NewAdminService(),
		Warehouse:        NewWarehouseService(repos.Warehouse),
		Service:          NewServiceService(repos.Service),
		ServiceOrderRepo: repository.NewServiceOrderRepository(db),
		Favorite:         NewFavoriteService(&repos.Favorite),
	}
}
