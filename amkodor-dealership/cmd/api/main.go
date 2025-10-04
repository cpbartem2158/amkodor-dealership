package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"amkodor-dealership/internal/config"
	"amkodor-dealership/internal/database"
	"amkodor-dealership/internal/handlers"
	"amkodor-dealership/internal/middleware"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/service"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

type Application struct {
	Config   *config.Config
	DB       *sql.DB
	Handlers *handlers.Handlers
}

func main() {
	// Загрузка .env файла
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Подключение к базе данных
	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Проверка соединения
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to PostgreSQL database")

	// Инициализация слоёв приложения
	app := initializeApplication(cfg, db)

	// Настройка роутера
	router := setupRouter(app)

	// Настройка CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	// Запуск сервера
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	srv := &http.Server{
		Handler:      corsHandler.Handler(router),
		Addr:         addr,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func initializeApplication(cfg *config.Config, db *sql.DB) *Application {
	// Инициализация репозиториев
	vehicleRepo := repository.NewVehicleRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	serviceRepo := repository.NewServiceRepository(db)

	// Инициализация сервисов
	vehicleService := service.NewVehicleService(vehicleRepo)
	customerService := service.NewCustomerService(customerRepo)
	saleService := service.NewSaleService(saleRepo, vehicleRepo)
	employeeService := service.NewEmployeeService(employeeRepo)
	authService := service.NewAuthService(employeeRepo, cfg.JWT.Secret)
	reportService := service.NewReportService(db)
	exportService := service.NewExportService(db)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	serviceOrderService := service.NewServiceOrderService(serviceRepo)

	// Инициализация обработчиков
	handlers := &handlers.Handlers{
		Vehicle:   handlers.NewVehicleHandler(vehicleService),
		Customer:  handlers.NewCustomerHandler(customerService),
		Sale:      handlers.NewSaleHandler(saleService),
		Employee:  handlers.NewEmployeeHandler(employeeService),
		Auth:      handlers.NewAuthHandler(authService),
		Dashboard: handlers.NewDashboardHandler(db),
		Report:    handlers.NewReportHandler(reportService, exportService),
		Admin:     handlers.NewAdminHandler(db),
		Warehouse: handlers.NewWarehouseHandler(warehouseService),
		Service:   handlers.NewServiceHandler(serviceOrderService),
	}

	return &Application{
		Config:   cfg,
		DB:       db,
		Handlers: handlers,
	}
}

func setupRouter(app *Application) *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(middleware.RecoveryMiddleware)

	// Статические файлы
	staticDir := "./web/static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Публичные маршруты
	r.HandleFunc("/", serveTemplate("index.html")).Methods("GET")
	r.HandleFunc("/login", serveTemplate("login.html")).Methods("GET")

	// API - Публичные эндпоинты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/auth/login", app.Handlers.Auth.Login).Methods("POST")
	api.HandleFunc("/vehicles", app.Handlers.Vehicle.GetAll).Methods("GET")
	api.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.GetByID).Methods("GET")
	api.HandleFunc("/vehicles/search", app.Handlers.Vehicle.Search).Methods("GET")
	api.HandleFunc("/test-drives", app.Handlers.Service.CreateTestDrive).Methods("POST")

	// API - Защищенные эндпоинты (требуют JWT)
	protected := api.PathPrefix("/admin").Subrouter()
	protected.Use(middleware.AuthMiddleware(app.Config.JWT.Secret))

	// Dashboard
	protected.HandleFunc("/dashboard", app.Handlers.Dashboard.GetStatistics).Methods("GET")
	protected.HandleFunc("/dashboard/charts", app.Handlers.Dashboard.GetChartData).Methods("GET")

	// Vehicles - CRUD
	protected.HandleFunc("/vehicles", app.Handlers.Vehicle.Create).Methods("POST")
	protected.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.Update).Methods("PUT")
	protected.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.Delete).Methods("DELETE")
	protected.HandleFunc("/vehicles/{id}/history", app.Handlers.Vehicle.GetHistory).Methods("GET")

	// Customers - CRUD
	protected.HandleFunc("/customers", app.Handlers.Customer.GetAll).Methods("GET")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.GetByID).Methods("GET")
	protected.HandleFunc("/customers", app.Handlers.Customer.Create).Methods("POST")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.Update).Methods("PUT")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.Delete).Methods("DELETE")
	protected.HandleFunc("/customers/search", app.Handlers.Customer.Search).Methods("GET")

	// Corporate Clients
	protected.HandleFunc("/corporate-clients", app.Handlers.Customer.GetAllCorporate).Methods("GET")
	protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.GetCorporateByID).Methods("GET")
	protected.HandleFunc("/corporate-clients", app.Handlers.Customer.CreateCorporate).Methods("POST")
	protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.UpdateCorporate).Methods("PUT")
	protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.DeleteCorporate).Methods("DELETE")

	// Sales - CRUD
	protected.HandleFunc("/sales", app.Handlers.Sale.GetAll).Methods("GET")
	protected.HandleFunc("/sales/{id}", app.Handlers.Sale.GetByID).Methods("GET")
	protected.HandleFunc("/sales", app.Handlers.Sale.Create).Methods("POST")
	protected.HandleFunc("/sales/{id}", app.Handlers.Sale.Update).Methods("PUT")
	protected.HandleFunc("/sales/{id}", app.Handlers.Sale.Delete).Methods("DELETE")
	protected.HandleFunc("/sales/{id}/history", app.Handlers.Sale.GetHistory).Methods("GET")

	// Employees - CRUD
	protected.HandleFunc("/employees", app.Handlers.Employee.GetAll).Methods("GET")
	protected.HandleFunc("/employees/{id}", app.Handlers.Employee.GetByID).Methods("GET")
	protected.HandleFunc("/employees", app.Handlers.Employee.Create).Methods("POST")
	protected.HandleFunc("/employees/{id}", app.Handlers.Employee.Update).Methods("PUT")
	protected.HandleFunc("/employees/{id}", app.Handlers.Employee.Delete).Methods("DELETE")

	// Warehouses
	protected.HandleFunc("/warehouses", app.Handlers.Warehouse.GetAll).Methods("GET")
	protected.HandleFunc("/warehouses/{id}", app.Handlers.Warehouse.GetByID).Methods("GET")
	protected.HandleFunc("/warehouses", app.Handlers.Warehouse.Create).Methods("POST")
	protected.HandleFunc("/warehouses/{id}", app.Handlers.Warehouse.Update).Methods("PUT")
	protected.HandleFunc("/warehouses/{id}/statistics", app.Handlers.Warehouse.GetStatistics).Methods("GET")

	// Service Orders
	protected.HandleFunc("/service-orders", app.Handlers.Service.GetAllOrders).Methods("GET")
	protected.HandleFunc("/service-orders/{id}", app.Handlers.Service.GetOrderByID).Methods("GET")
	protected.HandleFunc("/service-orders", app.Handlers.Service.CreateOrder).Methods("POST")
	protected.HandleFunc("/service-orders/{id}", app.Handlers.Service.UpdateOrder).Methods("PUT")
	protected.HandleFunc("/service-orders/{id}/complete", app.Handlers.Service.CompleteOrder).Methods("POST")

	// Test Drives
	protected.HandleFunc("/test-drives", app.Handlers.Service.GetAllTestDrives).Methods("GET")
	protected.HandleFunc("/test-drives/{id}", app.Handlers.Service.UpdateTestDrive).Methods("PUT")

	// Spare Parts
	protected.HandleFunc("/spare-parts", app.Handlers.Service.GetAllParts).Methods("GET")
	protected.HandleFunc("/spare-parts/{id}", app.Handlers.Service.GetPartByID).Methods("GET")
	protected.HandleFunc("/spare-parts", app.Handlers.Service.CreatePart).Methods("POST")
	protected.HandleFunc("/spare-parts/{id}", app.Handlers.Service.UpdatePart).Methods("PUT")

	// Reports
	protected.HandleFunc("/reports/sales", app.Handlers.Report.SalesReport).Methods("GET")
	protected.HandleFunc("/reports/inventory", app.Handlers.Report.InventoryReport).Methods("GET")
	protected.HandleFunc("/reports/export/sales", app.Handlers.Report.ExportSalesReport).Methods("GET")
	protected.HandleFunc("/reports/export/inventory", app.Handlers.Report.ExportInventoryReport).Methods("GET")

	// Admin Panel - Template routes
	adminPanel := r.PathPrefix("/admin").Subrouter()
	adminPanel.Use(middleware.AuthMiddleware(app.Config.JWT.Secret))
	adminPanel.HandleFunc("/dashboard", serveTemplate("admin/dashboard.html")).Methods("GET")
	adminPanel.HandleFunc("/vehicles", serveTemplate("admin/vehicles.html")).Methods("GET")
	adminPanel.HandleFunc("/sales", serveTemplate("admin/sales.html")).Methods("GET")
	adminPanel.HandleFunc("/customers", serveTemplate("admin/customers.html")).Methods("GET")
	adminPanel.HandleFunc("/employees", serveTemplate("admin/employees.html")).Methods("GET")
	adminPanel.HandleFunc("/service", serveTemplate("admin/service.html")).Methods("GET")
	adminPanel.HandleFunc("/reports", serveTemplate("admin/reports.html")).Methods("GET")
	adminPanel.HandleFunc("/settings", serveTemplate("admin/settings.html")).Methods("GET")

	return r
}

func serveTemplate(templatePath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullPath := "./web/templates/" + templatePath
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.Error(w, "Page not found", http.StatusNotFound)
			return
		}
		http.ServeFile(w, r, fullPath)
	}
}
