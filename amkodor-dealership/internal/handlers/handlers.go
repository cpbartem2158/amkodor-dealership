package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"

	"github.com/gorilla/mux"
)

// Структура для группировки всех handlers
type Handlers struct {
	Vehicle   *VehicleHandler
	Customer  *CustomerHandler
	Sale      *SaleHandler
	Employee  *EmployeeHandler
	Auth      *AuthHandler
	Dashboard *DashboardHandler
	Report    *ReportHandler
	Admin     *AdminHandler
	Warehouse *WarehouseHandler
	Service   *ServiceHandler
}

// VehicleHandler
type VehicleHandler struct {
	service *service.VehicleService
}

func NewVehicleHandler(service *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{service: service}
}

func (h *VehicleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	vehicles, err := h.service.GetAll(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching vehicles")
		return
	}

	utils.RespondSuccess(w, vehicles)
}

func (h *VehicleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	vehicle, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Vehicle not found")
		return
	}

	utils.RespondSuccess(w, vehicle)
}

func (h *VehicleHandler) Search(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{
		"model_name":        r.URL.Query().Get("model_name"),
		"category_name":     r.URL.Query().Get("category_name"),
		"type_name":         r.URL.Query().Get("type_name"),
		"manufacturer_name": r.URL.Query().Get("manufacturer_name"),
		"status":            r.URL.Query().Get("status"),
		"city":              r.URL.Query().Get("city"),
	}

	vehicles, err := h.service.Search(params)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error searching vehicles")
		return
	}

	utils.RespondSuccess(w, vehicles)
}

func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.Create(&vehicle)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating vehicle")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"vehicle_id": id})
}

func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	vehicle.VehicleID = id
	if err := h.service.Update(&vehicle); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating vehicle")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Vehicle updated successfully")
}

func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting vehicle")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Vehicle deleted successfully")
}

func (h *VehicleHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	history, err := h.service.GetHistory(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching history")
		return
	}

	utils.RespondSuccess(w, history)
}

// CustomerHandler
type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	customers, err := h.service.GetAll(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching customers")
		return
	}

	utils.RespondSuccess(w, customers)
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	customer, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Customer not found")
		return
	}

	utils.RespondSuccess(w, customer)
}

func (h *CustomerHandler) Search(w http.ResponseWriter, r *http.Request) {
	params := map[string]interface{}{
		"search_term":  r.URL.Query().Get("search_term"),
		"phone":        r.URL.Query().Get("phone"),
		"email":        r.URL.Query().Get("email"),
		"min_discount": r.URL.Query().Get("min_discount"),
	}

	customers, err := h.service.Search(params)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error searching customers")
		return
	}

	utils.RespondSuccess(w, customers)
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.Create(&customer)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating customer")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"customer_id": id})
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var customer models.Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	customer.CustomerID = id
	if err := h.service.Update(&customer); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating customer")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Customer updated successfully")
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting customer")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Customer deleted successfully")
}

func (h *CustomerHandler) GetAllCorporate(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	clients, err := h.service.GetAllCorporate(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching corporate clients")
		return
	}

	utils.RespondSuccess(w, clients)
}

func (h *CustomerHandler) GetCorporateByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	client, err := h.service.GetCorporateByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Corporate client not found")
		return
	}

	utils.RespondSuccess(w, client)
}

func (h *CustomerHandler) CreateCorporate(w http.ResponseWriter, r *http.Request) {
	var client models.CorporateClient
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.CreateCorporate(&client)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating corporate client")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"corporate_client_id": id})
}

func (h *CustomerHandler) UpdateCorporate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var client models.CorporateClient
	if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	client.CorporateClientID = id
	if err := h.service.UpdateCorporate(&client); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating corporate client")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Corporate client updated successfully")
}

func (h *CustomerHandler) DeleteCorporate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.service.DeleteCorporate(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting corporate client")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Corporate client deleted successfully")
}

// EmployeeHandler
type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	employees, err := h.service.GetAll(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching employees")
		return
	}

	utils.RespondSuccess(w, employees)
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	employee, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Employee not found")
		return
	}

	utils.RespondSuccess(w, employee)
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		models.Employee
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.Create(&req.Employee, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating employee")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"employee_id": id})
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var employee models.Employee
	if err := json.NewDecoder(r.Body).Decode(&employee); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	employee.EmployeeID = id
	if err := h.service.Update(&employee); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating employee")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Employee updated successfully")
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error deleting employee")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Employee deleted successfully")
}

// AuthHandler
type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	token, employee, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"token": token,
		"user": map[string]interface{}{
			"id":    employee.EmployeeID,
			"email": employee.Email,
			"name":  employee.FirstName + " " + employee.LastName,
		},
	})
}

// DashboardHandler
type DashboardHandler struct {
	db *sql.DB
}

func NewDashboardHandler(db *sql.DB) *DashboardHandler {
	return &DashboardHandler{db: db}
}

func (h *DashboardHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	query := `SELECT * FROM vw_dashboard_statistics`

	var stats models.DashboardStatistics
	err := h.db.QueryRow(query).Scan(
		&stats.AvailableVehicles, &stats.SalesLastMonth, &stats.RevenueLastMonth,
		&stats.TotalCustomers, &stats.TotalCorporateClients,
		&stats.UpcomingTestDrives, &stats.ActiveServiceOrders,
	)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching statistics")
		return
	}

	utils.RespondSuccess(w, stats)
}

func (h *DashboardHandler) GetChartData(w http.ResponseWriter, r *http.Request) {
	// Получение данных для графиков
	chartData := map[string]interface{}{
		"months":         []string{"Янв", "Фев", "Мар", "Апр", "Май", "Июн"},
		"sales":          []int{12, 19, 15, 25, 22, 30},
		"categories":     []string{"Дорожно-строительная", "Погрузочная", "Коммунальная"},
		"category_sales": []int{45, 30, 25},
	}

	utils.RespondSuccess(w, chartData)
}

// ReportHandler
type ReportHandler struct {
	reportService *service.ReportService
	exportService *service.ExportService
}

func NewReportHandler(reportService *service.ReportService, exportService *service.ExportService) *ReportHandler {
	return &ReportHandler{
		reportService: reportService,
		exportService: exportService,
	}
}

func (h *ReportHandler) SalesReport(w http.ResponseWriter, r *http.Request) {
	startDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", r.URL.Query().Get("end_date"))

	sales, err := h.reportService.GenerateSalesReport(startDate, endDate, nil, nil, nil)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error generating report")
		return
	}

	utils.RespondSuccess(w, sales)
}

func (h *ReportHandler) InventoryReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.reportService.GenerateInventoryReport(nil, nil, nil)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error generating report")
		return
	}

	utils.RespondSuccess(w, report)
}

func (h *ReportHandler) ExportSalesReport(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Excel export
	utils.RespondMessage(w, http.StatusOK, "Export functionality coming soon")
}

func (h *ReportHandler) ExportInventoryReport(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement Excel export
	utils.RespondMessage(w, http.StatusOK, "Export functionality coming soon")
}

// AdminHandler
type AdminHandler struct {
	db *sql.DB
}

func NewAdminHandler(db *sql.DB) *AdminHandler {
	return &AdminHandler{db: db}
}

// WarehouseHandler
type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(service *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching warehouses")
		return
	}

	utils.RespondSuccess(w, warehouses)
}

func (h *WarehouseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	warehouse, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Warehouse not found")
		return
	}

	utils.RespondSuccess(w, warehouse)
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var warehouse models.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.Create(&warehouse)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating warehouse")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"warehouse_id": id})
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var warehouse models.Warehouse
	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	warehouse.WarehouseID = id
	if err := h.service.Update(&warehouse); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating warehouse")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Warehouse updated successfully")
}

func (h *WarehouseHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	stats, err := h.service.GetStatistics(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching statistics")
		return
	}

	utils.RespondSuccess(w, stats)
}

// ServiceHandler
type ServiceHandler struct {
	service *service.ServiceOrderService
}

func NewServiceHandler(service *service.ServiceOrderService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

func (h *ServiceHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	orders, err := h.service.GetAllOrders(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching orders")
		return
	}

	utils.RespondSuccess(w, orders)
}

func (h *ServiceHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	order, err := h.service.GetOrderByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Order not found")
		return
	}

	utils.RespondSuccess(w, order)
}

func (h *ServiceHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var order models.ServiceOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.CreateOrder(&order)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating order")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"service_order_id": id})
}

func (h *ServiceHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var order models.ServiceOrder
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	order.ServiceOrderID = id
	if err := h.service.UpdateOrder(&order); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating order")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Order updated successfully")
}

func (h *ServiceHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	if err := h.service.CompleteOrder(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error completing order")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Order completed successfully")
}

func (h *ServiceHandler) GetAllParts(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	parts, err := h.service.GetAllParts(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching parts")
		return
	}

	utils.RespondSuccess(w, parts)
}

func (h *ServiceHandler) GetPartByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	part, err := h.service.GetPartByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Part not found")
		return
	}

	utils.RespondSuccess(w, part)
}

func (h *ServiceHandler) CreatePart(w http.ResponseWriter, r *http.Request) {
	var part models.SparePart
	if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.CreatePart(&part)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating part")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"spare_part_id": id})
}

func (h *ServiceHandler) UpdatePart(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var part models.SparePart
	if err := json.NewDecoder(r.Body).Decode(&part); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	part.SparePartID = id
	if err := h.service.UpdatePart(&part); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating part")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Part updated successfully")
}

func (h *ServiceHandler) GetAllTestDrives(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit == 0 {
		limit = 50
	}

	testDrives, err := h.service.GetAllTestDrives(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error fetching test drives")
		return
	}

	utils.RespondSuccess(w, testDrives)
}

func (h *ServiceHandler) CreateTestDrive(w http.ResponseWriter, r *http.Request) {
	var testDrive models.TestDrive
	if err := json.NewDecoder(r.Body).Decode(&testDrive); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	id, err := h.service.CreateTestDrive(&testDrive)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error creating test drive")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{"test_drive_id": id})
}

func (h *ServiceHandler) UpdateTestDrive(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var testDrive models.TestDrive
	if err := json.NewDecoder(r.Body).Decode(&testDrive); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid input")
		return
	}

	testDrive.TestDriveID = id
	if err := h.service.UpdateTestDrive(&testDrive); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Error updating test drive")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Test drive updated successfully")
}
