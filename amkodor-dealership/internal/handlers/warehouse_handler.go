package handlers

import (
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type WarehouseHandler struct {
	service *service.WarehouseService
}

func NewWarehouseHandler(service *service.WarehouseService) *WarehouseHandler {
	return &WarehouseHandler{service: service}
}

func (h *WarehouseHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	warehouses, err := h.service.GetAll()
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка получения списка складов")
		return
	}

	utils.RespondSuccess(w, warehouses)
}

func (h *WarehouseHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID склада")
		return
	}

	warehouse, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Склад не найден")
		return
	}

	utils.RespondSuccess(w, warehouse)
}

func (h *WarehouseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var warehouse models.Warehouse

	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	// Валидация обязательных полей
	if warehouse.WarehouseName == "" || warehouse.Address == "" || warehouse.City == "" {
		utils.RespondError(w, http.StatusBadRequest, "Заполните все обязательные поля")
		return
	}

	id, err := h.service.Create(&warehouse)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка создания склада")
		return
	}

	warehouse.WarehouseID = id
	utils.RespondSuccess(w, warehouse)
}

func (h *WarehouseHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID склада")
		return
	}

	var warehouse models.Warehouse

	if err := json.NewDecoder(r.Body).Decode(&warehouse); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный формат данных")
		return
	}

	// Валидация обязательных полей
	if warehouse.WarehouseName == "" || warehouse.Address == "" || warehouse.City == "" {
		utils.RespondError(w, http.StatusBadRequest, "Заполните все обязательные поля")
		return
	}

	warehouse.WarehouseID = id
	err = h.service.Update(&warehouse)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка обновления склада")
		return
	}

	utils.RespondSuccess(w, warehouse)
}

func (h *WarehouseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID склада")
		return
	}

	err = h.service.Delete(id)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка удаления склада")
		return
	}

	utils.RespondSuccess(w, map[string]string{"message": "Склад успешно удален"})
}

func (h *WarehouseHandler) GetStatistics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Неверный ID склада")
		return
	}

	// Получаем склад с дополнительной информацией
	warehouse, err := h.service.GetByID(id)
	if err != nil {
		utils.RespondError(w, http.StatusNotFound, "Склад не найден")
		return
	}

	utils.RespondSuccess(w, warehouse)
}
