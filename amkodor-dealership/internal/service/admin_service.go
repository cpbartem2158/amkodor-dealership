package service

import (
	"amkodor-dealership/internal/models"
)

type AdminService struct {
	// Заглушка
}

func NewAdminService() *AdminService {
	return &AdminService{}
}

func (s *AdminService) GetUsers() ([]models.User, error) {
	return []models.User{}, nil
}

func (s *AdminService) CreateUser(user *models.User) (int, error) {
	return 1, nil
}

func (s *AdminService) UpdateUser(user *models.User) error {
	return nil
}

func (s *AdminService) DeleteUser(id int) error {
	return nil
}
