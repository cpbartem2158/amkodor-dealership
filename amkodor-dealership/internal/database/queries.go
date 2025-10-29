package database

// SQL запросы для работы с базой данных

const (
	// Vehicles queries
	GetAllVehiclesQuery = `
		SELECT v.id, v.vin, v.model, v.year, v.color, v.mileage, v.price, v.status, 
		       v.category_id, v.created_at, v.updated_at,
		       c.name as category_name
		FROM vehicles v
		LEFT JOIN categories c ON v.category_id = c.id
		ORDER BY v.created_at DESC
	`

	GetVehicleByIDQuery = `
		SELECT v.id, v.vin, v.model, v.year, v.color, v.mileage, v.price, v.status,
		       v.category_id, v.created_at, v.updated_at,
		       c.name as category_name, c.description as category_description
		FROM vehicles v
		LEFT JOIN categories c ON v.category_id = c.id
		WHERE v.id = $1
	`

	CreateVehicleQuery = `
		INSERT INTO vehicles (vin, model, year, color, mileage, price, status, category_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	UpdateVehicleQuery = `
		UPDATE vehicles
		SET vin = $2, model = $3, year = $4, color = $5, mileage = $6,
		    price = $7, status = $8, category_id = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	DeleteVehicleQuery = `DELETE FROM vehicles WHERE id = $1`

	SearchVehiclesQuery = `
		SELECT v.id, v.vin, v.model, v.year, v.color, v.mileage, v.price, v.status,
		       v.category_id, v.created_at, c.name as category_name
		FROM vehicles v
		LEFT JOIN categories c ON v.category_id = c.id
		WHERE ($1::text IS NULL OR v.model ILIKE '%' || $1 || '%')
		  AND ($2::int IS NULL OR v.year >= $2)
		  AND ($3::int IS NULL OR v.year <= $3)
		  AND ($4::numeric IS NULL OR v.price >= $4)
		  AND ($5::numeric IS NULL OR v.price <= $5)
		  AND ($6::text IS NULL OR v.status = $6)
		  AND ($7::text IS NULL OR v.color ILIKE '%' || $7 || '%')
		ORDER BY v.created_at DESC
	`

	// Categories queries
	GetAllCategoriesQuery = `
		SELECT id, name, description, created_at
		FROM categories
		ORDER BY name
	`

	// Customers queries
	GetAllCustomersQuery = `
		SELECT c.id, c.full_name, c.phone, c.email, c.address, c.passport_number, c.created_at,
		       COUNT(s.id) as purchases_count,
		       COALESCE(SUM(s.final_price), 0) as total_spent
		FROM customers c
		LEFT JOIN sales s ON c.id = s.customer_id AND s.status = 'completed'
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	GetCustomerByIDQuery = `
		SELECT id, full_name, phone, email, address, passport_number, created_at, updated_at
		FROM customers
		WHERE id = $1
	`

	CreateCustomerQuery = `
		INSERT INTO customers (full_name, phone, email, address, passport_number)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`

	UpdateCustomerQuery = `
		UPDATE customers
		SET full_name = $2, phone = $3, email = $4, address = $5,
		    passport_number = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	DeleteCustomerQuery = `DELETE FROM customers WHERE id = $1`

	SearchCustomersQuery = `
		SELECT c.id, c.full_name, c.phone, c.email, c.address, c.passport_number, c.created_at,
		       COUNT(s.id) as purchases_count,
		       COALESCE(SUM(s.final_price), 0) as total_spent
		FROM customers c
		LEFT JOIN sales s ON c.id = s.customer_id
		WHERE ($1::text IS NULL OR c.full_name ILIKE '%' || $1 || '%')
		  AND ($2::text IS NULL OR c.phone LIKE '%' || $2 || '%')
		  AND ($3::text IS NULL OR c.email ILIKE '%' || $3 || '%')
		GROUP BY c.id
		ORDER BY c.created_at DESC
	`

	GetCustomerPurchaseHistoryQuery = `
		SELECT s.id, s.sale_date, s.discount, s.final_price, s.status,
		       v.model, v.vin, v.year
		FROM sales s
		JOIN vehicles v ON s.vehicle_id = v.id
		WHERE s.customer_id = $1
		ORDER BY s.sale_date DESC
	`

	// Sales queries
	GetAllSalesQuery = `
		SELECT s.id, s.sale_date, s.discount, s.final_price, s.status, s.created_at,
		       v.id as vehicle_id, v.model, v.vin, v.year,
		       c.id as customer_id, c.full_name as customer_name,
		       e.id as employee_id, e.full_name as employee_name
		FROM sales s
		JOIN vehicles v ON s.vehicle_id = v.id
		JOIN customers c ON s.customer_id = c.id
		JOIN employees e ON s.employee_id = e.id
		ORDER BY s.sale_date DESC
	`

	GetSaleByIDQuery = `
		SELECT s.id, s.sale_date, s.discount, s.final_price, s.status, s.created_at, s.updated_at,
		       v.id as vehicle_id, v.model, v.vin, v.year, v.price,
		       c.id as customer_id, c.full_name as customer_name, c.phone, c.email,
		       e.id as employee_id, e.full_name as employee_name
		FROM sales s
		JOIN vehicles v ON s.vehicle_id = v.id
		JOIN customers c ON s.customer_id = c.id
		JOIN employees e ON s.employee_id = e.id
		WHERE s.id = $1
	`

	CreateSaleQuery = `
		INSERT INTO sales (vehicle_id, customer_id, employee_id, sale_date, discount, final_price, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	UpdateSaleQuery = `
		UPDATE sales
		SET discount = $2, final_price = $3, status = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	CancelSaleQuery = `
		UPDATE sales
		SET status = 'cancelled', updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	// Employees queries
	GetAllEmployeesQuery = `
		SELECT id, full_name, position, phone, email, hire_date, salary, commission_rate, created_at
		FROM employees
		ORDER BY hire_date DESC
	`

	GetEmployeeByIDQuery = `
		SELECT id, full_name, position, phone, email, hire_date, salary, commission_rate, created_at, updated_at
		FROM employees
		WHERE id = $1
	`

	CreateEmployeeQuery = `
		INSERT INTO employees (full_name, position, phone, email, hire_date, salary, commission_rate)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`

	UpdateEmployeeQuery = `
		UPDATE employees
		SET full_name = $2, position = $3, phone = $4, email = $5,
		    salary = $6, commission_rate = $7, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	DeleteEmployeeQuery = `DELETE FROM employees WHERE id = $1`

	GetEmployeeSalesStatsQuery = `
		SELECT 
		    COUNT(*) as total_sales,
		    COALESCE(SUM(final_price), 0) as total_revenue,
		    COALESCE(AVG(final_price), 0) as avg_sale_price
		FROM sales
		WHERE employee_id = $1 AND status = 'completed'
	`

	// Dashboard queries
	GetDashboardStatsQuery = `
		SELECT 
		    (SELECT COUNT(*) FROM vehicles WHERE status = 'available') as available_vehicles,
		    (SELECT COUNT(*) FROM customers) as total_customers,
		    (SELECT COUNT(*) FROM sales WHERE status = 'completed') as total_sales,
		    (SELECT COALESCE(SUM(final_price), 0) FROM sales WHERE status = 'completed') as total_revenue,
		    (SELECT COUNT(*) FROM employees) as total_employees,
		    (SELECT COUNT(*) FROM sales WHERE sale_date >= CURRENT_DATE - INTERVAL '30 days' AND status = 'completed') as sales_last_month
	`

	GetSalesChartDataQuery = `
		SELECT 
		    DATE(sale_date) as date,
		    COUNT(*) as count,
		    SUM(final_price) as revenue
		FROM sales
		WHERE sale_date >= $1 AND status = 'completed'
		GROUP BY DATE(sale_date)
		ORDER BY date
	`

	GetTopEmployeesQuery = `
		SELECT 
		    e.id, e.full_name, e.position,
		    COUNT(s.id) as sales_count,
		    COALESCE(SUM(s.final_price), 0) as total_revenue
		FROM employees e
		LEFT JOIN sales s ON e.id = s.employee_id AND s.status = 'completed'
		GROUP BY e.id, e.full_name, e.position
		ORDER BY total_revenue DESC
		LIMIT $1
	`

	GetRecentSalesQuery = `
		SELECT s.id, s.sale_date, s.final_price, s.status,
		       v.model, c.full_name as customer_name, e.full_name as employee_name
		FROM sales s
		JOIN vehicles v ON s.vehicle_id = v.id
		JOIN customers c ON s.customer_id = c.id
		JOIN employees e ON s.employee_id = e.id
		ORDER BY s.sale_date DESC
		LIMIT $1
	`

	// Reports queries
	GenerateSalesReportQuery = `
		SELECT 
		    s.id, s.sale_date, s.discount, s.final_price, s.status,
		    v.model, v.vin, v.year, v.price as vehicle_price,
		    c.full_name as customer_name, c.phone as customer_phone,
		    e.full_name as employee_name,
		    (s.final_price - v.price) as profit
		FROM sales s
		JOIN vehicles v ON s.vehicle_id = v.id
		JOIN customers c ON s.customer_id = c.id
		JOIN employees e ON s.employee_id = e.id
		WHERE s.sale_date BETWEEN $1 AND $2
		  AND ($3::int IS NULL OR s.employee_id = $3)
		  AND s.status = 'completed'
		ORDER BY s.sale_date DESC
	`

	GenerateEmployeesReportQuery = `
		SELECT 
		    e.id, e.full_name, e.position, e.commission_rate,
		    COUNT(s.id) as sales_count,
		    COALESCE(SUM(s.final_price), 0) as total_revenue,
		    COALESCE(SUM(s.final_price) * e.commission_rate / 100, 0) as total_commission
		FROM employees e
		LEFT JOIN sales s ON e.id = s.employee_id 
		    AND s.sale_date BETWEEN $1 AND $2 
		    AND s.status = 'completed'
		GROUP BY e.id, e.full_name, e.position, e.commission_rate
		ORDER BY total_revenue DESC
	`

	GenerateVehiclesReportQuery = `
		SELECT 
		    c.name as category,
		    COUNT(v.id) as total_count,
		    SUM(CASE WHEN v.status = 'available' THEN 1 ELSE 0 END) as available_count,
		    SUM(CASE WHEN v.status = 'sold' THEN 1 ELSE 0 END) as sold_count,
		    AVG(v.price) as avg_price
		FROM vehicles v
		JOIN categories c ON v.category_id = c.id
		GROUP BY c.name
		ORDER BY c.name
	`

	// Auth queries
	GetUserByEmailQuery = `
		SELECT id, full_name, email, password_hash, position, created_at
		FROM employees
		WHERE email = $1
	`

	// Count queries
	CountVehiclesQuery  = `SELECT COUNT(*) FROM vehicles`
	CountCustomersQuery = `SELECT COUNT(*) FROM customers`
	CountEmployeesQuery = `SELECT COUNT(*) FROM employees`
	CountSalesQuery     = `SELECT COUNT(*) FROM sales WHERE status = 'completed'`
)
