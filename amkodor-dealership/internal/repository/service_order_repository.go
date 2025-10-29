package repository

import (
	"amkodor-dealership/internal/models"
	"database/sql"
	"fmt"
	"log"
)

type ServiceOrderRepository struct {
	db *sql.DB
}

func NewServiceOrderRepository(db *sql.DB) *ServiceOrderRepository {
	return &ServiceOrderRepository{db: db}
}

// CreateServiceOrder создает новый сервисный заказ
func (r *ServiceOrderRepository) CreateServiceOrder(req models.CreateServiceOrderRequest) (*models.ServiceOrder, error) {
	query := `
		INSERT INTO service_orders (vehicle_id, customer_id, corporate_client_id, employee_id, service_type, description, cost, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING service_order_id, order_date, created_at`

	var order models.ServiceOrder
	var customerID, corporateClientID sql.NullInt64
	var description sql.NullString

	if req.CustomerID != nil {
		customerID = sql.NullInt64{Int64: int64(*req.CustomerID), Valid: true}
	}
	if req.CorporateClientID != nil {
		corporateClientID = sql.NullInt64{Int64: int64(*req.CorporateClientID), Valid: true}
	}
	if req.Description != "" {
		description = sql.NullString{String: req.Description, Valid: true}
	}

	err := r.db.QueryRow(query, req.VehicleID, customerID, corporateClientID, req.EmployeeID, req.ServiceType, description, req.Cost, req.Status).
		Scan(&order.ServiceOrderID, &order.OrderDate, &order.CreatedAt)

	if err != nil {
		log.Printf("ERROR creating service order: %v, vehicle_id=%d, employee_id=%d", err, req.VehicleID, req.EmployeeID)
		return nil, fmt.Errorf("failed to create service order: %w", err)
	}

	order.VehicleID = req.VehicleID
	order.CustomerID = customerID
	order.CorporateClientID = corporateClientID
	order.EmployeeID = req.EmployeeID
	order.ServiceType = req.ServiceType
	order.Description = description
	order.Cost = req.Cost
	order.Status = req.Status

	return &order, nil
}

// GetAllServiceOrders получает все сервисные заказы
func (r *ServiceOrderRepository) GetAllServiceOrders() ([]models.ServiceOrder, error) {
	query := `
		SELECT service_order_id, vehicle_id, customer_id, corporate_client_id, employee_id, 
		       order_date, completion_date, service_type, description, cost, status, created_at
		FROM service_orders
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get service orders: %w", err)
	}
	defer rows.Close()

	var orders []models.ServiceOrder
	for rows.Next() {
		var order models.ServiceOrder
		err := rows.Scan(&order.ServiceOrderID, &order.VehicleID, &order.CustomerID, &order.CorporateClientID,
			&order.EmployeeID, &order.OrderDate, &order.CompletionDate, &order.ServiceType,
			&order.Description, &order.Cost, &order.Status, &order.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service order: %w", err)
		}
		orders = append(orders, order)
	}

	return orders, nil
}

// CreateTestDrive создает новый тест-драйв
func (r *ServiceOrderRepository) CreateTestDrive(req models.CreateTestDriveRequest) (*models.TestDrive, error) {
	query := `
		INSERT INTO test_drives (vehicle_id, customer_id, corporate_client_id, employee_id, scheduled_date, duration, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING test_drive_id, created_at`

	var testDrive models.TestDrive
	var customerID, corporateClientID sql.NullInt64

	if req.CustomerID != nil {
		customerID = sql.NullInt64{Int64: int64(*req.CustomerID), Valid: true}
	}
	if req.CorporateClientID != nil {
		corporateClientID = sql.NullInt64{Int64: int64(*req.CorporateClientID), Valid: true}
	}

	err := r.db.QueryRow(query, req.VehicleID, customerID, corporateClientID, req.EmployeeID, req.ScheduledDate, req.Duration, req.Status).
		Scan(&testDrive.TestDriveID, &testDrive.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create test drive: %w", err)
	}

	testDrive.VehicleID = req.VehicleID
	testDrive.CustomerID = customerID
	testDrive.CorporateClientID = corporateClientID
	testDrive.EmployeeID = req.EmployeeID
	testDrive.ScheduledDate = req.ScheduledDate
	testDrive.Duration = req.Duration
	testDrive.Status = req.Status

	return &testDrive, nil
}

// GetAllTestDrives получает все тест-драйвы
func (r *ServiceOrderRepository) GetAllTestDrives() ([]models.TestDrive, error) {
	query := `
		SELECT test_drive_id, vehicle_id, customer_id, corporate_client_id, employee_id, 
		       scheduled_date, duration, status, feedback_rating, feedback_comment, created_at
		FROM test_drives
		ORDER BY scheduled_date DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get test drives: %w", err)
	}
	defer rows.Close()

	var testDrives []models.TestDrive
	for rows.Next() {
		var testDrive models.TestDrive
		err := rows.Scan(&testDrive.TestDriveID, &testDrive.VehicleID, &testDrive.CustomerID, &testDrive.CorporateClientID,
			&testDrive.EmployeeID, &testDrive.ScheduledDate, &testDrive.Duration, &testDrive.Status,
			&testDrive.FeedbackRating, &testDrive.FeedbackComment, &testDrive.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan test drive: %w", err)
		}
		testDrives = append(testDrives, testDrive)
	}

	return testDrives, nil
}

// CreateSparePart создает новую запчасть
func (r *ServiceOrderRepository) CreateSparePart(req models.CreateSparePartRequest) (*models.SparePart, error) {
	query := `
		INSERT INTO spare_parts (part_number, part_name, model_id, price, quantity_in_stock, min_quantity, warehouse_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING spare_part_id, created_at`

	var part models.SparePart
	var modelID sql.NullInt64

	if req.ModelID != nil {
		modelID = sql.NullInt64{Int64: int64(*req.ModelID), Valid: true}
	}

	err := r.db.QueryRow(query, req.PartNumber, req.PartName, modelID, req.Price, req.QuantityInStock, req.MinQuantity, req.WarehouseID).
		Scan(&part.SparePartID, &part.CreatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create spare part: %w", err)
	}

	part.PartNumber = req.PartNumber
	part.PartName = req.PartName
	part.ModelID = modelID
	part.Price = req.Price
	part.QuantityInStock = req.QuantityInStock
	part.MinQuantity = req.MinQuantity
	part.WarehouseID = req.WarehouseID

	return &part, nil
}

// GetAllSpareParts получает все запчасти
func (r *ServiceOrderRepository) GetAllSpareParts() ([]models.SparePart, error) {
	query := `
		SELECT spare_part_id, part_number, part_name, model_id, price, quantity_in_stock, min_quantity, warehouse_id, created_at
		FROM spare_parts
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get spare parts: %w", err)
	}
	defer rows.Close()

	var parts []models.SparePart
	for rows.Next() {
		var part models.SparePart
		err := rows.Scan(&part.SparePartID, &part.PartNumber, &part.PartName, &part.ModelID, &part.Price,
			&part.QuantityInStock, &part.MinQuantity, &part.WarehouseID, &part.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan spare part: %w", err)
		}
		parts = append(parts, part)
	}

	return parts, nil
}

// DeleteSparePart удаляет запчасть
func (r *ServiceOrderRepository) DeleteSparePart(id int) error {
	query := `DELETE FROM spare_parts WHERE spare_part_id = $1`
	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete spare part: %w", err)
	}
	return nil
}
