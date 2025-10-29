package handlers

import (
	"amkodor-dealership/internal/service"
	"net/http"
)

type EmployeeHandler struct {
	service *service.EmployeeService
}

func NewEmployeeHandler(service *service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Employee handler not implemented"))
}

func (h *EmployeeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Employee handler not implemented"))
}

func (h *EmployeeHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Employee handler not implemented"))
}

func (h *EmployeeHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Employee handler not implemented"))
}

func (h *EmployeeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Employee handler not implemented"))
}
