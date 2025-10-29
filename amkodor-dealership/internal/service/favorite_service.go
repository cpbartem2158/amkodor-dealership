package service

import (
	"context"
	"fmt"
	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/repository"
)

type FavoriteService struct {
	repo *repository.FavoriteRepository
}

func NewFavoriteService(repo *repository.FavoriteRepository) *FavoriteService {
	return &FavoriteService{repo: repo}
}

// AddToFavorites добавляет технику в избранное
func (s *FavoriteService) AddToFavorites(ctx context.Context, userID, vehicleID int) error {
	return s.repo.AddToFavorites(ctx, userID, vehicleID)
}

// RemoveFromFavorites удаляет технику из избранного
func (s *FavoriteService) RemoveFromFavorites(ctx context.Context, userID, vehicleID int) error {
	return s.repo.RemoveFromFavorites(ctx, userID, vehicleID)
}

// GetUserFavorites получает все избранные пользователя
func (s *FavoriteService) GetUserFavorites(ctx context.Context, userID int) ([]models.FavoriteWithVehicle, error) {
	return s.repo.GetUserFavorites(ctx, userID)
}

// IsFavorite проверяет, находится ли техника в избранном
func (s *FavoriteService) IsFavorite(ctx context.Context, userID, vehicleID int) (bool, error) {
	return s.repo.IsFavorite(ctx, userID, vehicleID)
}

// GetFavoriteCount получает количество избранных у пользователя
func (s *FavoriteService) GetFavoriteCount(ctx context.Context, userID int) (int, error) {
	return s.repo.GetFavoriteCount(ctx, userID)
}

// ToggleFavorite переключает статус избранного (добавляет или удаляет)
func (s *FavoriteService) ToggleFavorite(ctx context.Context, userID, vehicleID int) (bool, error) {
	// Проверяем, есть ли уже в избранном
	isFavorite, err := s.repo.IsFavorite(ctx, userID, vehicleID)
	if err != nil {
		return false, fmt.Errorf("error checking favorite status: %w", err)
	}
	
	if isFavorite {
		// Удаляем из избранного
		err = s.repo.RemoveFromFavorites(ctx, userID, vehicleID)
		if err != nil {
			return false, fmt.Errorf("error removing from favorites: %w", err)
		}
		return false, nil
	} else {
		// Добавляем в избранное
		err = s.repo.AddToFavorites(ctx, userID, vehicleID)
		if err != nil {
			return false, fmt.Errorf("error adding to favorites: %w", err)
		}
		return true, nil
	}
}
