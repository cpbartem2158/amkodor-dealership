package main

import (
	"database/sql"
	"encoding/json"
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
	"amkodor-dealership/internal/utils"

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
	app, userService := initializeApplication(cfg, db)

	// Настройка роутера
	router := setupRouter(app, userService, db)

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

func initializeApplication(cfg *config.Config, db *sql.DB) (*Application, *service.UserService) {
	// Инициализация репозиториев
	vehicleRepo := repository.NewVehicleRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	saleRepo := repository.NewSaleRepository(db)
	employeeRepo := repository.NewEmployeeRepository(db)
	warehouseRepo := repository.NewWarehouseRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	serviceOrderRepo := repository.NewServiceOrderRepository(db)
	userRepo := repository.NewUserRepository(db)

	// Инициализация сервисов
	vehicleService := service.NewVehicleService(&vehicleRepo)
	customerService := service.NewCustomerService(customerRepo)
	saleService := service.NewSaleService(saleRepo, vehicleRepo)
	employeeService := service.NewEmployeeService(employeeRepo)
	authService := service.NewAuthService(&userRepo, cfg.JWT.Secret)
	userService := service.NewUserService(userRepo)
	reportService := service.NewReportService(db)
	_ = service.NewExportService(db)
	warehouseService := service.NewWarehouseService(warehouseRepo)
	_ = service.NewServiceOrderService(&serviceRepo)

	// Инициализация обработчиков
	handlers := &handlers.Handlers{
		Vehicle:   handlers.NewVehicleHandler(vehicleService),
		Customer:  handlers.NewCustomerHandler(customerService),
		Sale:      handlers.NewSaleHandler(saleService),
		Employee:  handlers.NewEmployeeHandler(employeeService),
		Auth:      handlers.NewAuthHandler(authService),
		Dashboard: handlers.NewDashboardHandler(warehouseService),
		Report:    handlers.NewReportHandler(reportService),
		Admin:     handlers.NewAdminHandler(warehouseService),
		Warehouse: handlers.NewWarehouseHandler(warehouseService),
		Service:   handlers.NewServiceHandler(serviceOrderRepo),
	}

	return &Application{
		Config:   cfg,
		DB:       db,
		Handlers: handlers,
	}, userService
}

func setupRouter(app *Application, userService *service.UserService, db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.RecoveryMiddleware)

	// Статические файлы
	staticDir := "./web/static"
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

	// Публичные маршруты
	r.HandleFunc("/", serveTemplate("index.html")).Methods("GET")
	r.HandleFunc("/catalog", serveTemplate("catalog.html")).Methods("GET")
	r.HandleFunc("/orders", serveTemplate("orders.html")).Methods("GET")
	r.HandleFunc("/login", serveTemplate("login.html")).Methods("GET")
	r.HandleFunc("/register", serveTemplate("register.html")).Methods("GET")
	r.HandleFunc("/dashboard", serveTemplate("dashboard.html")).Methods("GET")
	r.HandleFunc("/profile", serveTemplate("profile.html")).Methods("GET")
	r.HandleFunc("/favorites", serveTemplate("favorites.html")).Methods("GET")
	r.HandleFunc("/service", serveTemplate("service.html")).Methods("GET")

	// Админ панель
	r.HandleFunc("/admin/dashboard", serveTemplate("admin/dashboard.html")).Methods("GET")
	r.HandleFunc("/admin/vehicles", serveTemplate("admin/vehicles.html")).Methods("GET")
	r.HandleFunc("/admin/sales", serveTemplate("admin/sales.html")).Methods("GET")
	r.HandleFunc("/admin/customers", serveTemplate("admin/customers.html")).Methods("GET")
	r.HandleFunc("/admin/employees", serveTemplate("admin/employees.html")).Methods("GET")
	r.HandleFunc("/admin/warehouses", serveTemplate("admin/warehouses.html")).Methods("GET")
	r.HandleFunc("/admin/service", serveTemplate("admin/service.html")).Methods("GET")
	r.HandleFunc("/admin/reports", serveTemplate("admin/reports.html")).Methods("GET")
	r.HandleFunc("/admin/settings", serveTemplate("admin/settings.html")).Methods("GET")

	// Тестовая страница для загрузки изображений
	r.HandleFunc("/test_upload.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/test_upload.html")
	}).Methods("GET")

	// API - Публичные эндпоинты
	api := r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok","message":"Server is running"}`))
	}).Methods("GET")
	api.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		// Получаем реальную статистику из базы данных
		stats, err := getRealStats(db)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Ошибка получения статистики",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data":    stats,
		})
	}).Methods("GET")
	api.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем данные авторизации
		var loginData struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Неверный формат данных",
			})
			return
		}

		// Используем реальную авторизацию через базу данных
		user, err := userService.Login(r.Context(), loginData.Email, loginData.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Неверный email или пароль",
			})
			return
		}

		// Генерируем JWT токен
		token, err := utils.GenerateJWT(user.UserID, user.Email, app.Config.JWT.Secret, app.Config.JWT.ExpireHours)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Ошибка генерации токена",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"token":   token,
				"message": "Login successful",
				"role":    user.Role,
				"name":    user.Name,
			},
		})
	}).Methods("POST")

	api.HandleFunc("/auth/register", func(w http.ResponseWriter, r *http.Request) {
		// Проверяем данные регистрации
		var registerData struct {
			Name     string `json:"name"`
			Email    string `json:"email"`
			Phone    string `json:"phone"`
			Password string `json:"password"`
		}

		if err := json.NewDecoder(r.Body).Decode(&registerData); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Неверный формат данных",
			})
			return
		}

		// Проверяем обязательные поля
		if registerData.Name == "" || registerData.Email == "" || registerData.Password == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Все поля обязательны для заполнения",
			})
			return
		}

		// Проверяем длину пароля
		if len(registerData.Password) < 6 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Пароль должен содержать минимум 6 символов",
			})
			return
		}

		// Регистрируем пользователя через сервис
		user, err := userService.Register(r.Context(), registerData.Name, registerData.Email, registerData.Phone, registerData.Password)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		// Регистрируем пользователя
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"message": "Регистрация успешна",
				"user": map[string]string{
					"name":  user.Name,
					"email": user.Email,
					"phone": user.Phone,
				},
			},
		})
	}).Methods("POST")

	// Получение профиля пользователя
	api.HandleFunc("/auth/profile", func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Токен авторизации не найден",
			})
			return
		}

		// Извлекаем токен
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Валидируем токен и получаем пользователя
		user, err := userService.GetUserByToken(r.Context(), tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Недействительный токен",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"user_id": user.UserID,
				"name":    user.Name,
				"email":   user.Email,
				"phone":   user.Phone,
				"role":    user.Role,
			},
		})
	}).Methods("GET")

	// Обновление профиля пользователя
	api.HandleFunc("/auth/profile", func(w http.ResponseWriter, r *http.Request) {
		// Получаем токен из заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Токен авторизации не найден",
			})
			return
		}

		// Извлекаем токен
		tokenString := authHeader
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		}

		// Валидируем токен и получаем пользователя
		user, err := userService.GetUserByToken(r.Context(), tokenString)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Недействительный токен",
			})
			return
		}

		var req struct {
			Name  string `json:"name"`
			Phone string `json:"phone"`
		}

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Неверный формат данных",
			})
			return
		}

		// Обновляем данные пользователя
		user.Name = req.Name
		user.Phone = req.Phone

		err = userService.UpdateProfile(r.Context(), user)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"success": false,
				"error":   "Ошибка обновления профиля",
			})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"message": "Профиль успешно обновлен",
				"user": map[string]interface{}{
					"name":  user.Name,
					"email": user.Email,
					"phone": user.Phone,
				},
			},
		})
	}).Methods("PUT")

	api.HandleFunc("/vehicles", app.Handlers.Vehicle.GetAll).Methods("GET")
	api.HandleFunc("/vehicles", app.Handlers.Vehicle.Create).Methods("POST")
	api.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.GetByID).Methods("GET")
	api.HandleFunc("/vehicles/search", app.Handlers.Vehicle.Search).Methods("GET")
	api.HandleFunc("/vehicles/upload-image", app.Handlers.Vehicle.UploadImage).Methods("POST")
	// api.HandleFunc("/test-drives", app.Handlers.Service.CreateTestDrive).Methods("POST")

	// API - Избранное (требует JWT)
	api.HandleFunc("/favorites", app.Handlers.Favorite.GetUserFavorites).Methods("GET")
	api.HandleFunc("/favorites/{id}/toggle", app.Handlers.Favorite.ToggleFavorite).Methods("POST")
	api.HandleFunc("/favorites/{id}/check", app.Handlers.Favorite.IsFavorite).Methods("GET")
	api.HandleFunc("/favorites/count", app.Handlers.Favorite.GetFavoriteCount).Methods("GET")

	// API - Пользовательские функции (требуют JWT)
	api.HandleFunc("/user/stats", app.Handlers.User.GetUserStats).Methods("GET")
	api.HandleFunc("/user/orders", app.Handlers.User.GetUserOrders).Methods("GET")
	api.HandleFunc("/user/favorites", app.Handlers.User.GetUserFavorites).Methods("GET")
	api.HandleFunc("/user/favorites", app.Handlers.User.AddToFavorites).Methods("POST")
	api.HandleFunc("/user/favorites", app.Handlers.User.RemoveFromFavorites).Methods("DELETE")
	api.HandleFunc("/user/service-requests", app.Handlers.User.GetServiceRequests).Methods("GET")
	api.HandleFunc("/user/service-requests", app.Handlers.User.CreateServiceRequest).Methods("POST")
	api.HandleFunc("/user/profile", app.Handlers.User.GetUserProfile).Methods("GET")
	api.HandleFunc("/user/profile", app.Handlers.User.UpdateUserProfile).Methods("PUT")

	// API - Защищенные эндпоинты (требуют JWT)
	protected := api.PathPrefix("/admin").Subrouter()
	protected.Use(middleware.AuthMiddleware(app.Config.JWT.Secret))

	// Dashboard
	protected.HandleFunc("/dashboard", app.Handlers.Dashboard.GetStats).Methods("GET")
	protected.HandleFunc("/dashboard/charts", app.Handlers.Dashboard.GetCharts).Methods("GET")

	// Vehicles - CRUD
	protected.HandleFunc("/vehicles", app.Handlers.Vehicle.Create).Methods("POST")
	protected.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.Update).Methods("PUT")
	protected.HandleFunc("/vehicles/{id}", app.Handlers.Vehicle.Delete).Methods("DELETE")
	// protected.HandleFunc("/vehicles/{id}/history", app.Handlers.Vehicle.GetHistory).Methods("GET")

	// Customers - CRUD
	protected.HandleFunc("/customers", app.Handlers.Customer.GetAll).Methods("GET")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.GetByID).Methods("GET")
	protected.HandleFunc("/customers", app.Handlers.Customer.Create).Methods("POST")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.Update).Methods("PUT")
	protected.HandleFunc("/customers/{id}", app.Handlers.Customer.Delete).Methods("DELETE")
	// protected.HandleFunc("/customers/search", app.Handlers.Customer.Search).Methods("GET")

	// Corporate Clients - заглушки
	// protected.HandleFunc("/corporate-clients", app.Handlers.Customer.GetAllCorporate).Methods("GET")
	// protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.GetCorporateByID).Methods("GET")
	// protected.HandleFunc("/corporate-clients", app.Handlers.Customer.CreateCorporate).Methods("POST")
	// protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.UpdateCorporate).Methods("PUT")
	// protected.HandleFunc("/corporate-clients/{id}", app.Handlers.Customer.DeleteCorporate).Methods("DELETE")

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

	// Service Requests - получение всех заявок для админа
	protected.HandleFunc("/service-requests", app.Handlers.User.GetAllServiceRequests).Methods("GET")

	// Test Drives
	protected.HandleFunc("/test-drives", app.Handlers.Service.GetAllTestDrives).Methods("GET")
	protected.HandleFunc("/test-drives", app.Handlers.Service.CreateTestDrive).Methods("POST")
	protected.HandleFunc("/test-drives/{id}", app.Handlers.Service.UpdateTestDrive).Methods("PUT")

	// Spare Parts
	protected.HandleFunc("/spare-parts", app.Handlers.Service.GetAllParts).Methods("GET")
	protected.HandleFunc("/spare-parts/{id}", app.Handlers.Service.GetPartByID).Methods("GET")
	protected.HandleFunc("/spare-parts", app.Handlers.Service.CreatePart).Methods("POST")
	protected.HandleFunc("/spare-parts/{id}", app.Handlers.Service.UpdatePart).Methods("PUT")
	protected.HandleFunc("/spare-parts/{id}", app.Handlers.Service.DeleteSparePart).Methods("DELETE")

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
	adminPanel.HandleFunc("/warehouses", serveTemplate("admin/warehouses.html")).Methods("GET")
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

// getRealStats получает реальную статистику из базы данных
func getRealStats(db *sql.DB) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Общее количество техники
	var totalVehicles int
	err := db.QueryRow("SELECT COUNT(*) FROM vehicles").Scan(&totalVehicles)
	if err != nil {
		return nil, fmt.Errorf("error getting total vehicles: %w", err)
	}
	stats["total_vehicles"] = totalVehicles

	// Количество доступной техники
	var availableVehicles int
	err = db.QueryRow("SELECT COUNT(*) FROM vehicles WHERE status = 'В наличии'").Scan(&availableVehicles)
	if err != nil {
		return nil, fmt.Errorf("error getting available vehicles: %w", err)
	}
	stats["available_vehicles"] = availableVehicles

	// Общее количество продаж
	var totalSales int
	err = db.QueryRow("SELECT COUNT(*) FROM sales").Scan(&totalSales)
	if err != nil {
		return nil, fmt.Errorf("error getting total sales: %w", err)
	}
	stats["total_sales"] = totalSales

	// Общая выручка
	var totalRevenue sql.NullFloat64
	err = db.QueryRow("SELECT COALESCE(SUM(final_price), 0) FROM sales").Scan(&totalRevenue)
	if err != nil {
		return nil, fmt.Errorf("error getting total revenue: %w", err)
	}
	if totalRevenue.Valid {
		stats["total_revenue"] = int(totalRevenue.Float64)
	} else {
		stats["total_revenue"] = 0
	}

	// Общее количество клиентов
	var totalCustomers int
	err = db.QueryRow("SELECT COUNT(*) FROM customers").Scan(&totalCustomers)
	if err != nil {
		return nil, fmt.Errorf("error getting total customers: %w", err)
	}
	stats["total_customers"] = totalCustomers

	return stats, nil
}

// getRealVehicles получает реальные данные техники из базы данных
func getRealVehicles(db *sql.DB) ([]map[string]interface{}, error) {
	// Возвращаем пустой массив, так как мы используем новое API
	return []map[string]interface{}{}, nil
}
