package handlers

import (
	"net/http"
	"amkodor-dealership/internal/service"
)

type AdminHandler struct {
	service *service.WarehouseService
}

func NewAdminHandler(service *service.WarehouseService) *AdminHandler {
	return &AdminHandler{service: service}
}

func (h *AdminHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Admin handler not implemented"))
}

func (h *AdminHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Admin handler not implemented"))
}

func (h *AdminHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Admin handler not implemented"))
}

func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Admin handler not implemented"))
}