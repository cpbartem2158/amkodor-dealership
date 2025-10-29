package handlers

import (
	"amkodor-dealership/internal/service"
	"database/sql"
)

// Структура для группировки всех handlers
type Handlers struct {
	Vehicle   *VehicleHandler
	Customer  *CustomerHandler
	Sale      *SaleHandler
	Employee  *EmployeeHandler
	Auth      *AuthHandler
	Dashboard *DashboardHandler
	Report    *ReportHandler
	Admin     *AdminHandler
	Warehouse *WarehouseHandler
	Service   *ServiceHandler
	Favorite  *FavoriteHandler
	User      *UserHandler
}

// NewHandlers создает новый экземпляр Handlers
func NewHandlers(db *sql.DB, services *service.Services) *Handlers {
	return &Handlers{
		Vehicle:   NewVehicleHandler(services.Vehicle),
		Customer:  NewCustomerHandler(services.Customer),
		Sale:      NewSaleHandler(services.Sale),
		Employee:  NewEmployeeHandler(services.Employee),
		Auth:      NewAuthHandler(services.Auth),
		Dashboard: NewDashboardHandler(services.Warehouse),
		Report:    NewReportHandler(services.Report),
		Admin:     NewAdminHandler(services.Warehouse),
		Warehouse: NewWarehouseHandler(services.Warehouse),
		Service:   NewServiceHandler(services.ServiceOrderRepo),
		Favorite:  NewFavoriteHandler(services.Favorite),
		User:      NewUserHandler(),
	}
}
