package service

import (
	"time"
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type DashboardService struct {
	repo repository.DashboardRepository
}

func NewDashboardService(repo repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetStats() (*models.DashboardStats, error) {
	return s.repo.GetDashboardStats(nil)
}

func (s *DashboardService) GetCharts() ([]models.ChartData, error) {
	return s.repo.GetSalesChartData(nil, time.Now().AddDate(0, -1, 0))
}

func (s *DashboardService) GetTopEmployees(limit int) ([]models.TopEmployee, error) {
	return s.repo.GetTopEmployees(nil, limit)
}
