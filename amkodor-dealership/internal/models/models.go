package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Vehicle представляет единицу техники
type Vehicle struct {
	VehicleID       int            `json:"vehicle_id"`
	ModelID         int            `json:"model_id"`
	WarehouseID     int            `json:"warehouse_id"`
	VIN             sql.NullString `json:"vin"`
	SerialNumber    string         `json:"serial_number"`
	ManufactureYear int            `json:"manufacture_year"`
	Color           sql.NullString `json:"color"`
	Price           float64        `json:"price"`
	Discount        float64        `json:"discount"`
	Status          string         `json:"status"`
	ArrivalDate     time.Time      `json:"arrival_date"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	// Дополнительные поля из JOIN
	ModelName        string          `json:"model_name,omitempty"`
	TypeName         string          `json:"type_name,omitempty"`
	CategoryName     string          `json:"category_name,omitempty"`
	ManufacturerName string          `json:"manufacturer_name,omitempty"`
	WarehouseName    string          `json:"warehouse_name,omitempty"`
	WarehouseCity    string          `json:"warehouse_city,omitempty"`
	FinalPrice       float64         `json:"final_price,omitempty"`
	Description      sql.NullString  `json:"description,omitempty"`
	Specifications   json.RawMessage `json:"specifications,omitempty"`
}

// VehicleModel представляет модель техники
type VehicleModel struct {
	ModelID        int             `json:"model_id"`
	ModelName      string          `json:"model_name"`
	TypeID         int             `json:"type_id"`
	ManufacturerID int             `json:"manufacturer_id"`
	Description    sql.NullString  `json:"description"`
	Specifications json.RawMessage `json:"specifications"`
	CreatedAt      time.Time       `json:"created_at"`
}

// Customer представляет клиента (физическое лицо)
type Customer struct {
	CustomerID      int            `json:"customer_id"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	MiddleName      sql.NullString `json:"middle_name"`
	Phone           string         `json:"phone"`
	Email           sql.NullString `json:"email"`
	PassportNumber  sql.NullString `json:"passport_number"`
	Address         sql.NullString `json:"address"`
	DateOfBirth     sql.NullTime   `json:"date_of_birth"`
	DiscountPercent float64        `json:"discount_percent"`
	IsVIP           bool           `json:"is_vip"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	// Дополнительные поля
	FullName       string  `json:"full_name,omitempty"`
	CustomerLevel  string  `json:"customer_level,omitempty"`
	TotalPurchases int     `json:"total_purchases,omitempty"`
	TotalSpent     float64 `json:"total_spent,omitempty"`
}

// CorporateClient представляет корпоративного клиента
type CorporateClient struct {
	CorporateClientID int            `json:"corporate_client_id"`
	CompanyName       string         `json:"company_name"`
	TaxID             string         `json:"tax_id"`
	LegalAddress      string         `json:"legal_address"`
	ContactPerson     sql.NullString `json:"contact_person"`
	Phone             string         `json:"phone"`
	Email             sql.NullString `json:"email"`
	BankAccount       sql.NullString `json:"bank_account"`
	BankName          sql.NullString `json:"bank_name"`
	DiscountPercent   float64        `json:"discount_percent"`
	ContractNumber    sql.NullString `json:"contract_number"`
	ContractDate      sql.NullTime   `json:"contract_date"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	// Дополнительные поля
	TotalPurchases int     `json:"total_purchases,omitempty"`
	TotalSpent     float64 `json:"total_spent,omitempty"`
}

// Sale представляет продажу
type Sale struct {
	SaleID            int            `json:"sale_id"`
	VehicleID         int            `json:"vehicle_id"`
	CustomerID        sql.NullInt64  `json:"customer_id"`
	CorporateClientID sql.NullInt64  `json:"corporate_client_id"`
	EmployeeID        int            `json:"employee_id"`
	SaleDate          time.Time      `json:"sale_date"`
	BasePrice         float64        `json:"base_price"`
	DiscountAmount    float64        `json:"discount_amount"`
	FinalPrice        float64        `json:"final_price"`
	PaymentType       string         `json:"payment_type"`
	Status            string         `json:"status"`
	ContractNumber    sql.NullString `json:"contract_number"`
	Notes             sql.NullString `json:"notes"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	// Дополнительные поля из JOIN
	VIN           string `json:"vin,omitempty"`
	ModelName     string `json:"model_name,omitempty"`
	TypeName      string `json:"type_name,omitempty"`
	ClientName    string `json:"client_name,omitempty"`
	ClientPhone   string `json:"client_phone,omitempty"`
	ClientType    string `json:"client_type,omitempty"`
	ManagerName   string `json:"manager_name,omitempty"`
	PositionName  string `json:"position_name,omitempty"`
	WarehouseName string `json:"warehouse_name,omitempty"`
	WarehouseCity string `json:"warehouse_city,omitempty"`
}

// Employee представляет сотрудника
type Employee struct {
	EmployeeID   int             `json:"employee_id"`
	FirstName    string          `json:"first_name"`
	LastName     string          `json:"last_name"`
	MiddleName   sql.NullString  `json:"middle_name"`
	PositionID   int             `json:"position_id"`
	WarehouseID  sql.NullInt64   `json:"warehouse_id"`
	Email        sql.NullString  `json:"email"`
	Phone        sql.NullString  `json:"phone"`
	PasswordHash string          `json:"-"` // Не экспортируем в JSON
	HireDate     time.Time       `json:"hire_date"`
	Salary       sql.NullFloat64 `json:"salary"`
	IsActive     bool            `json:"is_active"`
	CreatedAt    time.Time       `json:"created_at"`
	// Дополнительные поля
	FullName       string `json:"full_name,omitempty"`
	PositionName   string `json:"position_name,omitempty"`
	WarehouseName  string `json:"warehouse_name,omitempty"`
	WarehouseCity  string `json:"warehouse_city,omitempty"`
	YearsOfService int    `json:"years_of_service,omitempty"`
}

// Warehouse представляет склад/филиал
type Warehouse struct {
	WarehouseID   int            `json:"warehouse_id"`
	WarehouseName string         `json:"warehouse_name"`
	Address       string         `json:"address"`
	City          string         `json:"city"`
	Region        sql.NullString `json:"region"`
	Phone         sql.NullString `json:"phone"`
	ManagerName   sql.NullString `json:"manager_name"`
	Capacity      int            `json:"capacity"`
	CreatedAt     time.Time      `json:"created_at"`
	IsActive      bool           `json:"is_active"`
	// Статистика
	VehiclesInStock     int     `json:"vehicles_in_stock,omitempty"`
	AvailableVehicles   int     `json:"available_vehicles,omitempty"`
	EmployeesCount      int     `json:"employees_count,omitempty"`
	TotalInventoryValue float64 `json:"total_inventory_value,omitempty"`
}

// ServiceOrder представляет сервисный заказ
type ServiceOrder struct {
	ServiceOrderID    int            `json:"service_order_id"`
	VehicleID         int            `json:"vehicle_id"`
	CustomerID        sql.NullInt64  `json:"customer_id"`
	CorporateClientID sql.NullInt64  `json:"corporate_client_id"`
	EmployeeID        int            `json:"employee_id"`
	OrderDate         time.Time      `json:"order_date"`
	CompletionDate    sql.NullTime   `json:"completion_date"`
	ServiceType       string         `json:"service_type"`
	Description       sql.NullString `json:"description"`
	Cost              float64        `json:"cost"`
	Status            string         `json:"status"`
	CreatedAt         time.Time      `json:"created_at"`
	// Дополнительные поля
	ModelName   string  `json:"model_name,omitempty"`
	VIN         string  `json:"vin,omitempty"`
	ClientName  string  `json:"client_name,omitempty"`
	ClientPhone string  `json:"client_phone,omitempty"`
	MasterName  string  `json:"master_name,omitempty"`
	PartsCost   float64 `json:"parts_cost,omitempty"`
}

// TestDrive представляет тест-драйв
type TestDrive struct {
	TestDriveID       int            `json:"test_drive_id"`
	VehicleID         int            `json:"vehicle_id"`
	CustomerID        sql.NullInt64  `json:"customer_id"`
	CorporateClientID sql.NullInt64  `json:"corporate_client_id"`
	EmployeeID        int            `json:"employee_id"`
	ScheduledDate     time.Time      `json:"scheduled_date"`
	Duration          int            `json:"duration"`
	Status            string         `json:"status"`
	FeedbackRating    sql.NullInt64  `json:"feedback_rating"`
	FeedbackComment   sql.NullString `json:"feedback_comment"`
	CreatedAt         time.Time      `json:"created_at"`
	// Дополнительные поля
	ModelName       string `json:"model_name,omitempty"`
	TypeName        string `json:"type_name,omitempty"`
	Color           string `json:"color,omitempty"`
	ManufactureYear int    `json:"manufacture_year,omitempty"`
	ClientName      string `json:"client_name,omitempty"`
	ClientPhone     string `json:"client_phone,omitempty"`
	ManagerName     string `json:"manager_name,omitempty"`
	ManagerPhone    string `json:"manager_phone,omitempty"`
	WarehouseName   string `json:"warehouse_name,omitempty"`
	WarehouseCity   string `json:"warehouse_city,omitempty"`
}

// SparePart представляет запасную часть
type SparePart struct {
	SparePartID     int           `json:"spare_part_id"`
	PartNumber      string        `json:"part_number"`
	PartName        string        `json:"part_name"`
	ModelID         sql.NullInt64 `json:"model_id"`
	Price           float64       `json:"price"`
	QuantityInStock int           `json:"quantity_in_stock"`
	MinQuantity     int           `json:"min_quantity"`
	WarehouseID     int           `json:"warehouse_id"`
	CreatedAt       time.Time     `json:"created_at"`
	// Дополнительные поля
	ModelName     string `json:"model_name,omitempty"`
	StockStatus   string `json:"stock_status,omitempty"`
	WarehouseName string `json:"warehouse_name,omitempty"`
	WarehouseCity string `json:"warehouse_city,omitempty"`
}

// History структуры
type VehicleHistory struct {
	HistoryID       int             `json:"history_id"`
	VehicleID       int             `json:"vehicle_id"`
	OperationType   string          `json:"operation_type"`
	OperationDate   time.Time       `json:"operation_date"`
	OldValue        json.RawMessage `json:"old_value"`
	NewValue        json.RawMessage `json:"new_value"`
	Username        string          `json:"username"`
	Hostname        sql.NullString  `json:"hostname"`
	ApplicationName sql.NullString  `json:"application_name"`
}

type SaleHistory struct {
	HistoryID       int             `json:"history_id"`
	SaleID          int             `json:"sale_id"`
	OperationType   string          `json:"operation_type"`
	OperationDate   time.Time       `json:"operation_date"`
	OldValue        json.RawMessage `json:"old_value"`
	NewValue        json.RawMessage `json:"new_value"`
	Username        string          `json:"username"`
	Hostname        sql.NullString  `json:"hostname"`
	ApplicationName sql.NullString  `json:"application_name"`
}

// Statistics структуры для дашборда
type DashboardStatistics struct {
	AvailableVehicles     int     `json:"available_vehicles"`
	SalesLastMonth        int     `json:"sales_last_month"`
	RevenueLastMonth      float64 `json:"revenue_last_month"`
	TotalCustomers        int     `json:"total_customers"`
	TotalCorporateClients int     `json:"total_corporate_clients"`
	UpcomingTestDrives    int     `json:"upcoming_test_drives"`
	ActiveServiceOrders   int     `json:"active_service_orders"`
}

// Pagination для списков
type Pagination struct {
	Page       int `json:"page"`
	PerPage    int `json:"per_page"`
	TotalPages int `json:"total_pages"`
	TotalItems int `json:"total_items"`
}

// APIResponse общий ответ API
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}
