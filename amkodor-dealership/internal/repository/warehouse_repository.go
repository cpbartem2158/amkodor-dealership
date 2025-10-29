package repository

import (
	"database/sql"
	"fmt"

	"amkodor-dealership/internal/models"
)

type WarehouseRepository struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	return WarehouseRepository{db: db}
}

// GetAll возвращает все склады
func (r *WarehouseRepository) GetAll() ([]models.Warehouse, error) {
	query := `
		SELECT * FROM vw_warehouses_statistics
		ORDER BY warehouse_name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error querying warehouses: %w", err)
	}
	defer rows.Close()

	warehouses := []models.Warehouse{}
	for rows.Next() {
		var w models.Warehouse
		err := rows.Scan(
			&w.WarehouseID, &w.WarehouseName, &w.City, &w.Region, &w.Capacity,
			&w.VehiclesInStock, &w.AvailableVehicles, &w.EmployeesCount, &w.TotalInventoryValue,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning warehouse: %w", err)
		}
		warehouses = append(warehouses, w)
	}

	return warehouses, nil
}

// GetByID возвращает склад по ID
func (r *WarehouseRepository) GetByID(id int) (*models.Warehouse, error) {
	query := `
		SELECT 
			warehouse_id, warehouse_name, address, city, region,
			phone, manager_name, capacity, created_at, is_active
		FROM warehouses
		WHERE warehouse_id = $1
	`

	var w models.Warehouse
	err := r.db.QueryRow(query, id).Scan(
		&w.WarehouseID, &w.WarehouseName, &w.Address, &w.City, &w.Region,
		&w.Phone, &w.ManagerName, &w.Capacity, &w.CreatedAt, &w.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying warehouse: %w", err)
	}

	return &w, nil
}

// Create создает новый склад
func (r *WarehouseRepository) Create(w *models.Warehouse) (int, error) {
	query := `
		INSERT INTO warehouses (
			warehouse_name, address, city, region, phone, manager_name, capacity
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING warehouse_id
	`

	var warehouseID int
	err := r.db.QueryRow(
		query,
		w.WarehouseName, w.Address, w.City, w.Region, w.Phone, w.ManagerName, w.Capacity,
	).Scan(&warehouseID)

	if err != nil {
		return 0, fmt.Errorf("error creating warehouse: %w", err)
	}

	return warehouseID, nil
}

// Update обновляет склад
func (r *WarehouseRepository) Update(w *models.Warehouse) error {
	query := `
		UPDATE warehouses SET
			warehouse_name = $1,
			address = $2,
			city = $3,
			region = $4,
			phone = $5,
			manager_name = $6,
			capacity = $7,
			is_active = $8
		WHERE warehouse_id = $9
	`

	result, err := r.db.Exec(
		query,
		w.WarehouseName, w.Address, w.City, w.Region, w.Phone,
		w.ManagerName, w.Capacity, w.IsActive, w.WarehouseID,
	)

	if err != nil {
		return fmt.Errorf("error updating warehouse: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("warehouse not found")
	}

	return nil
}

// Delete удаляет склад
func (r *WarehouseRepository) Delete(id int) error {
	query := `DELETE FROM warehouses WHERE warehouse_id = $1`

	_, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting warehouse: %w", err)
	}

	return nil
}

// GetStatistics возвращает статистику склада
func (r *WarehouseRepository) GetStatistics(id int) (*models.Warehouse, error) {
	query := `
		SELECT * FROM vw_warehouses_statistics
		WHERE warehouse_id = $1
	`

	var w models.Warehouse
	err := r.db.QueryRow(query, id).Scan(
		&w.WarehouseID, &w.WarehouseName, &w.City, &w.Region, &w.Capacity,
		&w.VehiclesInStock, &w.AvailableVehicles, &w.EmployeesCount, &w.TotalInventoryValue,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("warehouse not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying warehouse statistics: %w", err)
	}

	return &w, nil
}
