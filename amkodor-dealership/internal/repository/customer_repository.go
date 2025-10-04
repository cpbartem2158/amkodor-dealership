package repository

import (
	"database/sql"
	"fmt"

	"amkodor-dealership/internal/models"
)

type CustomerRepository struct {
	db *sql.DB
}

func NewCustomerRepository(db *sql.DB) *CustomerRepository {
	return &CustomerRepository{db: db}
}

// GetAll возвращает всех клиентов
func (r *CustomerRepository) GetAll(limit, offset int) ([]models.Customer, error) {
	query := `
		SELECT 
			customer_id, first_name, last_name, middle_name, phone, email,
			passport_number, address, date_of_birth, discount_percent, is_vip,
			created_at, updated_at
		FROM customers
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying customers: %w", err)
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.CustomerID, &c.FirstName, &c.LastName, &c.MiddleName, &c.Phone,
			&c.Email, &c.PassportNumber, &c.Address, &c.DateOfBirth,
			&c.DiscountPercent, &c.IsVIP, &c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}
		customers = append(customers, c)
	}

	return customers, nil
}

// GetByID возвращает клиента по ID
func (r *CustomerRepository) GetByID(id int) (*models.Customer, error) {
	query := `
		SELECT 
			customer_id, first_name, last_name, middle_name, phone, email,
			passport_number, address, date_of_birth, discount_percent, is_vip,
			created_at, updated_at
		FROM customers
		WHERE customer_id = $1
	`

	var c models.Customer
	err := r.db.QueryRow(query, id).Scan(
		&c.CustomerID, &c.FirstName, &c.LastName, &c.MiddleName, &c.Phone,
		&c.Email, &c.PassportNumber, &c.Address, &c.DateOfBirth,
		&c.DiscountPercent, &c.IsVIP, &c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("customer not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying customer: %w", err)
	}

	return &c, nil
}

// Create создает нового клиента
func (r *CustomerRepository) Create(c *models.Customer) (int, error) {
	query := `
		INSERT INTO customers (
			first_name, last_name, middle_name, phone, email,
			passport_number, address, date_of_birth, discount_percent, is_vip
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING customer_id
	`

	var customerID int
	err := r.db.QueryRow(
		query,
		c.FirstName, c.LastName, c.MiddleName, c.Phone, c.Email,
		c.PassportNumber, c.Address, c.DateOfBirth, c.DiscountPercent, c.IsVIP,
	).Scan(&customerID)

	if err != nil {
		return 0, fmt.Errorf("error creating customer: %w", err)
	}

	return customerID, nil
}

// Update обновляет клиента
func (r *CustomerRepository) Update(c *models.Customer) error {
	query := `
		UPDATE customers SET
			first_name = $1,
			last_name = $2,
			middle_name = $3,
			phone = $4,
			email = $5,
			passport_number = $6,
			address = $7,
			date_of_birth = $8,
			discount_percent = $9,
			is_vip = $10,
			updated_at = CURRENT_TIMESTAMP
		WHERE customer_id = $11
	`

	result, err := r.db.Exec(
		query,
		c.FirstName, c.LastName, c.MiddleName, c.Phone, c.Email,
		c.PassportNumber, c.Address, c.DateOfBirth, c.DiscountPercent,
		c.IsVIP, c.CustomerID,
	)

	if err != nil {
		return fmt.Errorf("error updating customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Delete удаляет клиента
func (r *CustomerRepository) Delete(id int) error {
	query := `DELETE FROM customers WHERE customer_id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting customer: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Search ищет клиентов (используя процедуру)
func (r *CustomerRepository) Search(params map[string]interface{}) ([]models.Customer, error) {
	query := `
		SELECT * FROM sp_search_customers($1, $2, $3, $4, $5)
	`

	rows, err := r.db.Query(
		query,
		params["search_term"],
		params["phone"],
		params["email"],
		params["is_vip"],
		params["min_discount"],
	)

	if err != nil {
		return nil, fmt.Errorf("error searching customers: %w", err)
	}
	defer rows.Close()

	customers := []models.Customer{}
	for rows.Next() {
		var c models.Customer
		err := rows.Scan(
			&c.CustomerID, &c.FullName, &c.Phone, &c.Email, &c.Address,
			&c.DiscountPercent, &c.IsVIP, &c.CustomerLevel, &c.TotalPurchases,
			&c.TotalSpent, &c.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning customer: %w", err)
		}
		customers = append(customers, c)
	}

	return customers, nil
}

// Count возвращает общее количество клиентов
func (r *CustomerRepository) Count() (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM customers`
	err := r.db.QueryRow(query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error counting customers: %w", err)
	}
	return count, nil
}

// GetAllCorporate возвращает всех корпоративных клиентов
func (r *CustomerRepository) GetAllCorporate(limit, offset int) ([]models.CorporateClient, error) {
	query := `
		SELECT 
			corporate_client_id, company_name, tax_id, legal_address,
			contact_person, phone, email, bank_account, bank_name,
			discount_percent, contract_number, contract_date,
			created_at, updated_at
		FROM corporate_clients
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("error querying corporate clients: %w", err)
	}
	defer rows.Close()

	clients := []models.CorporateClient{}
	for rows.Next() {
		var c models.CorporateClient
		err := rows.Scan(
			&c.CorporateClientID, &c.CompanyName, &c.TaxID, &c.LegalAddress,
			&c.ContactPerson, &c.Phone, &c.Email, &c.BankAccount, &c.BankName,
			&c.DiscountPercent, &c.ContractNumber, &c.ContractDate,
			&c.CreatedAt, &c.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning corporate client: %w", err)
		}
		clients = append(clients, c)
	}

	return clients, nil
}

// GetCorporateByID возвращает корпоративного клиента по ID
func (r *CustomerRepository) GetCorporateByID(id int) (*models.CorporateClient, error) {
	query := `
		SELECT 
			corporate_client_id, company_name, tax_id, legal_address,
			contact_person, phone, email, bank_account, bank_name,
			discount_percent, contract_number, contract_date,
			created_at, updated_at
		FROM corporate_clients
		WHERE corporate_client_id = $1
	`

	var c models.CorporateClient
	err := r.db.QueryRow(query, id).Scan(
		&c.CorporateClientID, &c.CompanyName, &c.TaxID, &c.LegalAddress,
		&c.ContactPerson, &c.Phone, &c.Email, &c.BankAccount, &c.BankName,
		&c.DiscountPercent, &c.ContractNumber, &c.ContractDate,
		&c.CreatedAt, &c.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("corporate client not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error querying corporate client: %w", err)
	}

	return &c, nil
}

// CreateCorporate создает корпоративного клиента
func (r *CustomerRepository) CreateCorporate(c *models.CorporateClient) (int, error) {
	query := `
		INSERT INTO corporate_clients (
			company_name, tax_id, legal_address, contact_person, phone, email,
			bank_account, bank_name, discount_percent, contract_number, contract_date
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING corporate_client_id
	`

	var id int
	err := r.db.QueryRow(
		query,
		c.CompanyName, c.TaxID, c.LegalAddress, c.ContactPerson, c.Phone, c.Email,
		c.BankAccount, c.BankName, c.DiscountPercent, c.ContractNumber, c.ContractDate,
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("error creating corporate client: %w", err)
	}

	return id, nil
}

// UpdateCorporate обновляет корпоративного клиента
func (r *CustomerRepository) UpdateCorporate(c *models.CorporateClient) error {
	query := `
		UPDATE corporate_clients SET
			company_name = $1,
			tax_id = $2,
			legal_address = $3,
			contact_person = $4,
			phone = $5,
			email = $6,
			bank_account = $7,
			bank_name = $8,
			discount_percent = $9,
			contract_number = $10,
			contract_date = $11,
			updated_at = CURRENT_TIMESTAMP
		WHERE corporate_client_id = $12
	`

	result, err := r.db.Exec(
		query,
		c.CompanyName, c.TaxID, c.LegalAddress, c.ContactPerson, c.Phone, c.Email,
		c.BankAccount, c.BankName, c.DiscountPercent, c.ContractNumber, c.ContractDate,
		c.CorporateClientID,
	)

	if err != nil {
		return fmt.Errorf("error updating corporate client: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("corporate client not found")
	}

	return nil
}

// DeleteCorporate удаляет корпоративного клиента
func (r *CustomerRepository) DeleteCorporate(id int) error {
	query := `DELETE FROM corporate_clients WHERE corporate_client_id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error deleting corporate client: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected: %w", err)
	}

	if rows == 0 {
		return fmt.Errorf("corporate client not found")
	}

	return nil
}
