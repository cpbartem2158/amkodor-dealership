package handlers

import (
	"amkodor-dealership/internal/service"
	"net/http"
)

type CustomerHandler struct {
	service *service.CustomerService
}

func NewCustomerHandler(service *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{service: service}
}

func (h *CustomerHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Customer handler not implemented"))
}

func (h *CustomerHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Customer handler not implemented"))
}

func (h *CustomerHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Customer handler not implemented"))
}

func (h *CustomerHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Customer handler not implemented"))
}

func (h *CustomerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Customer handler not implemented"))
}
