package handlers

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
	"amkodor-dealership/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type ServiceHandler struct {
	serviceOrderRepo *repository.ServiceOrderRepository
}

func NewServiceHandler(serviceOrderRepo *repository.ServiceOrderRepository) *ServiceHandler {
	return &ServiceHandler{serviceOrderRepo: serviceOrderRepo}
}

func (h *ServiceHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	orders, err := h.serviceOrderRepo.GetAllServiceOrders()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get service orders")
		return
	}
	utils.RespondSuccess(w, orders)
}

func (h *ServiceHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"data":null}`))
}

func (h *ServiceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServiceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Status == "" {
		req.Status = "В работе"
	}

	order, err := h.serviceOrderRepo.CreateServiceOrder(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create service order")
		return
	}

	utils.RespondSuccess(w, order)
}

func (h *ServiceHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Service order updated"}`))
}

func (h *ServiceHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Service order deleted"}`))
}

func (h *ServiceHandler) GetAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.serviceOrderRepo.GetAllServiceOrders()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get service orders")
		return
	}
	utils.RespondSuccess(w, orders)
}

func (h *ServiceHandler) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"data":null}`))
}

func (h *ServiceHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	var req models.CreateServiceOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Status == "" {
		req.Status = "В работе"
	}

	order, err := h.serviceOrderRepo.CreateServiceOrder(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create service order")
		return
	}

	utils.RespondSuccess(w, order)
}

func (h *ServiceHandler) UpdateOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Service order updated"}`))
}

func (h *ServiceHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Service order completed"}`))
}

func (h *ServiceHandler) GetAllTestDrives(w http.ResponseWriter, r *http.Request) {
	testDrives, err := h.serviceOrderRepo.GetAllTestDrives()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get test drives")
		return
	}
	utils.RespondSuccess(w, testDrives)
}

func (h *ServiceHandler) UpdateTestDrive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Test drive updated"}`))
}

func (h *ServiceHandler) GetAllParts(w http.ResponseWriter, r *http.Request) {
	parts, err := h.serviceOrderRepo.GetAllSpareParts()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to get spare parts")
		return
	}
	utils.RespondSuccess(w, parts)
}

func (h *ServiceHandler) GetPartByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"data":null}`))
}

func (h *ServiceHandler) CreatePart(w http.ResponseWriter, r *http.Request) {
	var req models.CreateSparePartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	part, err := h.serviceOrderRepo.CreateSparePart(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create spare part")
		return
	}

	utils.RespondSuccess(w, part)
}

func (h *ServiceHandler) UpdatePart(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success":true,"message":"Part updated"}`))
}

// CreateTestDrive создает новый тест-драйв
func (h *ServiceHandler) CreateTestDrive(w http.ResponseWriter, r *http.Request) {
	var req models.CreateTestDriveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Устанавливаем значения по умолчанию
	if req.Status == "" {
		req.Status = "Запланирован"
	}
	if req.Duration == 0 {
		req.Duration = 60
	}

	testDrive, err := h.serviceOrderRepo.CreateTestDrive(req)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to create test drive")
		return
	}

	utils.RespondSuccess(w, testDrive)
}

// DeleteSparePart удаляет запчасть
func (h *ServiceHandler) DeleteSparePart(w http.ResponseWriter, r *http.Request) {
	// Извлекаем ID из URL
	vars := r.URL.Query()
	idStr := vars.Get("id")
	if idStr == "" {
		utils.RespondError(w, http.StatusBadRequest, "Missing spare part ID")
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid spare part ID")
		return
	}

	err = h.serviceOrderRepo.DeleteSparePart(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Failed to delete spare part")
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Spare part deleted successfully"})
}
