package repository

import (
	"database/sql"
	"fmt"

	"amkodor-dealership/internal/models"
)

type ServiceRepository struct {
	db *sql.DB
}

func NewServiceRepository(db *sql.DB) *ServiceRepository {
	return &ServiceRepository{db: db}
}

// GetAllOrders возвращает все сервисные заказы
func (r *ServiceRepository) GetAllOrders(limit, offset int) ([]models.ServiceOrder, error) {
	query := `
		SELECT * FROM vw_service_orders_full_info
		ORDER BY order_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying service orders: %w", err)
	}
	defer rows.Close()

	orders := []models.ServiceOrder{}
	for rows.Next() {
		var so models.ServiceOrder
		err := rows.Scan(
			&so.ServiceOrderID, &so.OrderDate, &so.CompletionDate, &so.ServiceType,
			&so.Status, &so.Cost, &so.ModelName, &so.VIN, &so.ClientName,
			&so.ClientPhone, &so.MasterName, &so.PartsCost, &so.Description,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning service order: %w", err)
		}
		orders = append(orders, so)
	}

	return orders, nil
}

// GetOrderByID возвращает сервисный заказ по ID
func (r *ServiceRepository) GetOrderByID(id int) (*models.ServiceOrder, error) {
	query := `
		SELECT * FROM vw_service_orders_full_info
		WHERE service_order_id = $1
	`

	var so models.ServiceOrder
	err := r.db.QueryRow(query, id).Scan(
		&so.ServiceOrderID, &so.OrderDate, &so.CompletionDate, &so.ServiceType,
		&so.Status, &so.Cost, &so.ModelName, &so.VIN, &so.ClientName,
		&so.ClientPhone, &so.MasterName, &so.PartsCost, &so.Description,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("service order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying service order: %w", err)
	}

	return &so, nil
}

// CreateOrder создает сервисный заказ
func (r *ServiceRepository) CreateOrder(so *models.ServiceOrder) (int, error) {
	query := `
		SELECT sp_create_service_order($1, $2, $3, $4, $5, $6, $7)
	`

	var orderID int
	err := r.db.QueryRow(
		query,
		so.VehicleID, so.CustomerID, so.CorporateClientID, so.EmployeeID,
		so.ServiceType, so.Description, so.Cost,
	).Scan(&orderID)

	if err != nil {
		return 0, fmt.Errorf("error creating service order: %w", err)
	}

	return orderID, nil
}

// UpdateOrder обновляет сервисный заказ
func (r *ServiceRepository) UpdateOrder(so *models.ServiceOrder) error {
	query := `
		UPDATE service_orders SET
			service_type = $1,
			description = $2,
			cost = $3,
			status = $4,
			completion_date = $5
		WHERE service_order_id = $6
	`

	result, err := r.db.Exec(
		query,
		so.ServiceType, so.Description, so.Cost, so.Status,
		so.CompletionDate, so.ServiceOrderID,
	)

	if err != nil {
		return fmt.Errorf("error updating service order: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("service order not found")
	}

	return nil
}

// CompleteOrder завершает сервисный заказ
func (r *ServiceRepository) CompleteOrder(orderID int) error {
	query := `SELECT sp_complete_service_order($1)`

	_, err := r.db.Exec(query, orderID)
	if err != nil {
		return fmt.Errorf("error completing service order: %w", err)
	}

	return nil
}

// GetAllParts возвращает все запчасти
func (r *ServiceRepository) GetAllParts(limit, offset int) ([]models.SparePart, error) {
	query := `
		SELECT * FROM vw_spare_parts_inventory
		ORDER BY part_name
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying spare parts: %w", err)
	}
	defer rows.Close()

	parts := []models.SparePart{}
	for rows.Next() {
		var sp models.SparePart
		err := rows.Scan(
			&sp.SparePartID, &sp.PartNumber, &sp.PartName, &sp.ModelID,
			&sp.Price, &sp.QuantityInStock, &sp.MinQuantity, &sp.StockStatus,
			&sp.WarehouseName, &sp.WarehouseCity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning spare part: %w", err)
		}
		parts = append(parts, sp)
	}

	return parts, nil
}

// GetPartByID возвращает запчасть по ID
func (r *ServiceRepository) GetPartByID(id int) (*models.SparePart, error) {
	query := `
		SELECT * FROM vw_spare_parts_inventory
		WHERE spare_part_id = $1
	`

	var sp models.SparePart
	err := r.db.QueryRow(query, id).Scan(
		&sp.SparePartID, &sp.PartNumber, &sp.PartName, &sp.ModelID,
		&sp.Price, &sp.QuantityInStock, &sp.MinQuantity, &sp.StockStatus,
		&sp.WarehouseName, &sp.WarehouseCity,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("spare part not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying spare part: %w", err)
	}

	return &sp, nil
}

// CreatePart создает запчасть
func (r *ServiceRepository) CreatePart(sp *models.SparePart) (int, error) {
	query := `
		INSERT INTO spare_parts (
			part_number, part_name, model_id, price, quantity_in_stock, min_quantity, warehouse_id
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING spare_part_id
	`

	var partID int
	err := r.db.QueryRow(
		query,
		sp.PartNumber, sp.PartName, sp.ModelID, sp.Price,
		sp.QuantityInStock, sp.MinQuantity, sp.WarehouseID,
	).Scan(&partID)

	if err != nil {
		return 0, fmt.Errorf("error creating spare part: %w", err)
	}

	return partID, nil
}

// UpdatePart обновляет запчасть
func (r *ServiceRepository) UpdatePart(sp *models.SparePart) error {
	query := `
		UPDATE spare_parts SET
			part_number = $1,
			part_name = $2,
			model_id = $3,
			price = $4,
			quantity_in_stock = $5,
			min_quantity = $6,
			warehouse_id = $7
		WHERE spare_part_id = $8
	`

	result, err := r.db.Exec(
		query,
		sp.PartNumber, sp.PartName, sp.ModelID, sp.Price,
		sp.QuantityInStock, sp.MinQuantity, sp.WarehouseID, sp.SparePartID,
	)

	if err != nil {
		return fmt.Errorf("error updating spare part: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("spare part not found")
	}

	return nil
}

// GetAllTestDrives возвращает все тест-драйвы
func (r *ServiceRepository) GetAllTestDrives(limit, offset int) ([]models.TestDrive, error) {
	query := `
		SELECT * FROM vw_test_drives_full_info
		ORDER BY scheduled_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying test drives: %w", err)
	}
	defer rows.Close()

	testDrives := []models.TestDrive{}
	for rows.Next() {
		var td models.TestDrive
		err := rows.Scan(
			&td.TestDriveID, &td.ScheduledDate, &td.Duration, &td.Status,
			&td.ModelName, &td.TypeName, &td.Color, &td.ManufactureYear,
			&td.ClientName, &td.ClientPhone, &td.ManagerName, &td.ManagerPhone,
			&td.FeedbackRating, &td.FeedbackComment, &td.WarehouseName, &td.WarehouseCity,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning test drive: %w", err)
		}
		testDrives = append(testDrives, td)
	}

	return testDrives, nil
}

// CreateTestDrive создает тест-драйв
func (r *ServiceRepository) CreateTestDrive(td *models.TestDrive) (int, error) {
	query := `
		SELECT sp_create_test_drive($1, $2, $3, $4, $5, $6)
	`

	var testDriveID int
	err := r.db.QueryRow(
		query,
		td.VehicleID, td.CustomerID, td.CorporateClientID,
		td.EmployeeID, td.ScheduledDate, td.Duration,
	).Scan(&testDriveID)

	if err != nil {
		return 0, fmt.Errorf("error creating test drive: %w", err)
	}

	return testDriveID, nil
}

// UpdateTestDrive обновляет тест-драйв
func (r *ServiceRepository) UpdateTestDrive(td *models.TestDrive) error {
	query := `
		UPDATE test_drives SET
			status = $1,
			feedback_rating = $2,
			feedback_comment = $3
		WHERE test_drive_id = $4
	`

	result, err := r.db.Exec(
		query,
		td.Status, td.FeedbackRating, td.FeedbackComment, td.TestDriveID,
	)

	if err != nil {
		return fmt.Errorf("error updating test drive: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("test drive not found")
	}

	return nil
}
