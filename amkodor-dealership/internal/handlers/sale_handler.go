package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"

	"github.com/gorilla/mux"
)

type SaleHandler struct {
	service *service.SaleService
}

func NewSaleHandler(service *service.SaleService) *SaleHandler {
	return &SaleHandler{service: service}
}

// GetAll возвращает список всех продаж
func (h *SaleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Параметры пагинации
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 50
	offset := 0

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	sales, err := h.service.GetAll(limit, offset)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка получения списка продаж")
		return
	}

	utils.RespondSuccess(w, sales)
}

// GetByID возвращает продажу по ID
func (h *SaleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID")
		return
	}

	sale, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Продажа не найдена")
		return
	}

	utils.RespondSuccess(w, sale)
}

// Create создает новую продажу
func (h *SaleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID          int     `json:"vehicle_id"`
		CustomerID         *int    `json:"customer_id"`
		CorporateClientID  *int    `json:"corporate_client_id"`
		EmployeeID         int     `json:"employee_id"`
		PaymentType        string  `json:"payment_type"`
		AdditionalDiscount float64 `json:"additional_discount"`
		ContractNumber     string  `json:"contract_number"`
		Notes              string  `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	// Валидация
	if req.VehicleID == 0 || req.EmployeeID == 0 {
		utils.RespondError(w, http.StatusBadRequest, "Необходимо указать технику и менеджера")
		return
	}

	if req.CustomerID == nil && req.CorporateClientID == nil {
		utils.RespondError(w, http.StatusBadRequest, "Необходимо указать клиента")
		return
	}

	if req.PaymentType == "" {
		req.PaymentType = "Наличные"
	}

	saleID, err := h.service.Create(
		req.VehicleID,
		req.CustomerID,
		req.CorporateClientID,
		req.EmployeeID,
		req.PaymentType,
		req.AdditionalDiscount,
		req.ContractNumber,
		req.Notes,
	)

	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"sale_id": saleID,
		"message": "Продажа успешно создана",
	})
}

// Update обновляет продажу
func (h *SaleHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID")
		return
	}

	var sale models.Sale
	if err := json.NewDecoder(r.Body).Decode(&sale); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	sale.SaleID = id

	if err := h.service.Update(&sale); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка обновления продажи")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Продажа успешно обновлена")
}

// Delete удаляет продажу
func (h *SaleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка удаления продажи")
		return
	}

	utils.RespondMessage(w, http.StatusOK, "Продажа успешно удалена")
}

// GetHistory возвращает историю изменений продажи
func (h *SaleHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID")
		return
	}

	history, err := h.service.GetHistory(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка получения истории")
		return
	}

	utils.RespondSuccess(w, history)
}
