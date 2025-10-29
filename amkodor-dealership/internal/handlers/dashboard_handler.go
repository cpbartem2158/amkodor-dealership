package handlers

import (
	"net/http"
	"amkodor-dealership/internal/service"
)

type DashboardHandler struct {
	service *service.WarehouseService
}

func NewDashboardHandler(service *service.WarehouseService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Dashboard handler not implemented"))
}

func (h *DashboardHandler) GetCharts(w http.ResponseWriter, r *http.Request) {
	// Заглушка
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("Dashboard handler not implemented"))
}