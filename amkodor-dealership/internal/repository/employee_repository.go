package repository

import (
	"database/sql"
	"fmt"

	"amkodor-dealership/internal/models"
)

type EmployeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) *EmployeeRepository {
	return &EmployeeRepository{db: db}
}

// GetAll возвращает всех сотрудников
func (r *EmployeeRepository) GetAll(limit, offset int) ([]models.Employee, error) {
	query := `
		SELECT * FROM vw_employees_full_info
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying employees: %w", err)
	}
	defer rows.Close()

	employees := []models.Employee{}
	for rows.Next() {
		var e models.Employee
		err := rows.Scan(
			&e.EmployeeID, &e.FullName, &e.FirstName, &e.LastName, &e.MiddleName,
			&e.PositionName, &e.Salary, &e.WarehouseName, &e.WarehouseCity,
			&e.Email, &e.Phone, &e.HireDate, &e.YearsOfService, &e.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning employee: %w", err)
		}
		employees = append(employees, e)
	}

	return employees, nil
}

// GetByID возвращает сотрудника по ID
func (r *EmployeeRepository) GetByID(id int) (*models.Employee, error) {
	query := `
		SELECT * FROM vw_employees_full_info
		WHERE employee_id = $1
	`

	var e models.Employee
	err := r.db.QueryRow(query, id).Scan(
		&e.EmployeeID, &e.FullName, &e.FirstName, &e.LastName, &e.MiddleName,
		&e.PositionName, &e.Salary, &e.WarehouseName, &e.WarehouseCity,
		&e.Email, &e.Phone, &e.HireDate, &e.YearsOfService, &e.IsActive,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying employee: %w", err)
	}

	return &e, nil
}

// GetByEmail возвращает сотрудника по email (для аутентификации)
func (r *EmployeeRepository) GetByEmail(email string) (*models.Employee, error) {
	query := `
		SELECT 
			employee_id, first_name, last_name, middle_name, position_id,
			warehouse_id, email, phone, password_hash, hire_date, salary, is_active, created_at
		FROM employees
		WHERE email = $1
	`

	var e models.Employee
	err := r.db.QueryRow(query, email).Scan(
		&e.EmployeeID, &e.FirstName, &e.LastName, &e.MiddleName, &e.PositionID,
		&e.WarehouseID, &e.Email, &e.Phone, &e.PasswordHash, &e.HireDate,
		&e.Salary, &e.IsActive, &e.CreatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("employee not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying employee: %w", err)
	}

	return &e, nil
}

// Create создает нового сотрудника
func (r *EmployeeRepository) Create(e *models.Employee) (int, error) {
	query := `
		INSERT INTO employees (
			first_name, last_name, middle_name, position_id, warehouse_id,
			email, phone, password_hash, salary
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING employee_id
	`

	var employeeID int
	err := r.db.QueryRow(
		query,
		e.FirstName, e.LastName, e.MiddleName, e.PositionID, e.WarehouseID,
		e.Email, e.Phone, e.PasswordHash, e.Salary,
	).Scan(&employeeID)

	if err != nil {
		return 0, fmt.Errorf("error creating employee: %w", err)
	}

	return employeeID, nil
}

// Update обновляет сотрудника
func (r *EmployeeRepository) Update(e *models.Employee) error {
	query := `
		UPDATE employees SET
			first_name = $1,
			last_name = $2,
			middle_name = $3,
			position_id = $4,
			warehouse_id = $5,
			email = $6,
			phone = $7,
			salary = $8,
			is_active = $9
		WHERE employee_id = $10
	`

	result, err := r.db.Exec(
		query,
		e.FirstName, e.LastName, e.MiddleName, e.PositionID, e.WarehouseID,
		e.Email, e.Phone, e.Salary, e.IsActive, e.EmployeeID,
	)

	if err != nil {
		return fmt.Errorf("error updating employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// UpdatePassword обновляет пароль сотрудника
func (r *EmployeeRepository) UpdatePassword(employeeID int, passwordHash string) error {
	query := `
		UPDATE employees SET
			password_hash = $1
		WHERE employee_id = $2
	`

	result, err := r.db.Exec(query, passwordHash, employeeID)
	if err != nil {
		return fmt.Errorf("error updating password: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// Delete удаляет сотрудника
func (r *EmployeeRepository) Delete(id int) error {
	query := `DELETE FROM employees WHERE employee_id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting employee: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("employee not found")
	}

	return nil
}

// Count возвращает общее количество сотрудников
func (r *EmployeeRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM employees WHERE is_active = true`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting employees: %w", err)
	}
	return count, nil
}
