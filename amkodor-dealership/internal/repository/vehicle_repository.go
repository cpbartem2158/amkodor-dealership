package repository

import (
	"database/sql"
	"fmt"

	"amkodor-dealership/internal/models"
)

type VehicleRepository struct {
	db *sql.DB
}

func NewVehicleRepository(db *sql.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

// GetAll возвращает все единицы техники
func (r *VehicleRepository) GetAll(limit, offset int) ([]models.Vehicle, error) {
	query := `
		SELECT * FROM vw_vehicles_full_info
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying vehicles: %w", err)
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.VehicleID, &v.VIN, &v.SerialNumber, &v.ModelName, &v.TypeName,
			&v.CategoryName, &v.ManufacturerName, &v.ManufactureYear, &v.Color,
			&v.Price, &v.Discount, &v.FinalPrice, &v.Status, &v.WarehouseName,
			&v.WarehouseCity, &v.ArrivalDate, &v.CreatedAt, &v.Description,
			&v.Specifications,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// GetByID возвращает технику по ID
func (r *VehicleRepository) GetByID(id int) (*models.Vehicle, error) {
	query := `
		SELECT * FROM vw_vehicles_full_info
		WHERE vehicle_id = $1
	`

	var v models.Vehicle
	err := r.db.QueryRow(query, id).Scan(
		&v.VehicleID, &v.VIN, &v.SerialNumber, &v.ModelName, &v.TypeName,
		&v.CategoryName, &v.ManufacturerName, &v.ManufactureYear, &v.Color,
		&v.Price, &v.Discount, &v.FinalPrice, &v.Status, &v.WarehouseName,
		&v.WarehouseCity, &v.ArrivalDate, &v.CreatedAt, &v.Description,
		&v.Specifications,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("vehicle not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying vehicle: %w", err)
	}

	return &v, nil
}

// Create создает новую единицу техники
func (r *VehicleRepository) Create(v *models.Vehicle) (int, error) {
	query := `
		INSERT INTO vehicles (
			model_id, warehouse_id, vin, serial_number, manufacture_year,
			color, price, discount, status, arrival_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING vehicle_id
	`

	var vehicleID int
	err := r.db.QueryRow(
		query,
		v.ModelID, v.WarehouseID, v.VIN, v.SerialNumber, v.ManufactureYear,
		v.Color, v.Price, v.Discount, v.Status, v.ArrivalDate,
	).Scan(&vehicleID)

	if err != nil {
		return 0, fmt.Errorf("error creating vehicle: %w", err)
	}

	return vehicleID, nil
}

// Update обновляет технику
func (r *VehicleRepository) Update(v *models.Vehicle) error {
	query := `
		UPDATE vehicles SET
			model_id = $1,
			warehouse_id = $2,
			vin = $3,
			serial_number = $4,
			manufacture_year = $5,
			color = $6,
			price = $7,
			discount = $8,
			status = $9,
			updated_at = CURRENT_TIMESTAMP
		WHERE vehicle_id = $10
	`

	result, err := r.db.Exec(
		query,
		v.ModelID, v.WarehouseID, v.VIN, v.SerialNumber, v.ManufactureYear,
		v.Color, v.Price, v.Discount, v.Status, v.VehicleID,
	)

	if err != nil {
		return fmt.Errorf("error updating vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

// Delete удаляет технику
func (r *VehicleRepository) Delete(id int) error {
	query := `DELETE FROM vehicles WHERE vehicle_id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting vehicle: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}

// Search ищет технику по параметрам (используя хранимую процедуру)
func (r *VehicleRepository) Search(params map[string]interface{}) ([]models.Vehicle, error) {
	query := `
		SELECT * FROM sp_search_vehicles(
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`

	rows, err := r.db.Query(
		query,
		params["model_name"],
		params["category_name"],
		params["type_name"],
		params["manufacturer_name"],
		params["min_price"],
		params["max_price"],
		params["min_year"],
		params["max_year"],
		params["status"],
		params["warehouse_id"],
		params["city"],
	)

	if err != nil {
		return nil, fmt.Errorf("error searching vehicles: %w", err)
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		var v models.Vehicle
		err := rows.Scan(
			&v.VehicleID, &v.VIN, &v.SerialNumber, &v.ModelName, &v.TypeName,
			&v.CategoryName, &v.ManufacturerName, &v.ManufactureYear, &v.Color,
			&v.Price, &v.Discount, &v.FinalPrice, &v.Status, &v.WarehouseName,
			&v.WarehouseCity, &v.WarehouseCity, &v.Description, &v.Specifications,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// GetHistory возвращает историю изменений техники
func (r *VehicleRepository) GetHistory(vehicleID int) ([]models.VehicleHistory, error) {
	query := `
		SELECT 
			history_id, vehicle_id, operation_type, operation_date,
			old_value, new_value, username, hostname, application_name
		FROM vehicles_history
		WHERE vehicle_id = $1
		ORDER BY operation_date DESC
	`

	rows, err := r.db.Query(query, vehicleID)
	if err != nil {
		return nil, fmt.Errorf("error querying vehicle history: %w", err)
	}
	defer rows.Close()

	history := []models.VehicleHistory{}
	for rows.Next() {
		var h models.VehicleHistory
		err := rows.Scan(
			&h.HistoryID, &h.VehicleID, &h.OperationType, &h.OperationDate,
			&h.OldValue, &h.NewValue, &h.Username, &h.Hostname, &h.ApplicationName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle history: %w", err)
		}
		history = append(history, h)
	}

	return history, nil
}

// GetAvailable возвращает доступную для продажи технику
func (r *VehicleRepository) GetAvailable(limit, offset int) ([]models.Vehicle, error) {
	query := `
		SELECT * FROM vw_available_vehicles
		ORDER BY arrival_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying available vehicles: %w", err)
	}
	defer rows.Close()

	vehicles := []models.Vehicle{}
	for rows.Next() {
		var v models.Vehicle
		var daysInStock int
		err := rows.Scan(
			&v.VehicleID, &v.ModelName, &v.TypeName, &v.CategoryName,
			&v.ManufacturerName, &v.ManufactureYear, &v.Color, &v.Price,
			&v.Discount, &v.FinalPrice, &v.WarehouseName, &v.WarehouseCity,
			&v.WarehouseCity, &v.Description, &v.Specifications, &daysInStock,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning vehicle: %w", err)
		}
		vehicles = append(vehicles, v)
	}

	return vehicles, nil
}

// Count возвращает общее количество техники
func (r *VehicleRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM vehicles`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting vehicles: %w", err)
	}
	return count, nil
}

// UpdateStatus обновляет статус техники
func (r *VehicleRepository) UpdateStatus(vehicleID int, status string) error {
	query := `
		UPDATE vehicles 
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE vehicle_id = $2
	`

	result, err := r.db.Exec(query, status, vehicleID)
	if err != nil {
		return fmt.Errorf("error updating vehicle status: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("vehicle not found")
	}

	return nil
}
