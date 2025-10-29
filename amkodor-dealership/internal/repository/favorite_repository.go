package repository

import (
	"context"
	"database/sql"
	"fmt"
	"amkodor-dealership/internal/models"
)

type FavoriteRepository struct {
	db *sql.DB
}

func NewFavoriteRepository(db *sql.DB) FavoriteRepository {
	return FavoriteRepository{db: db}
}

// AddToFavorites добавляет технику в избранное
func (r *FavoriteRepository) AddToFavorites(ctx context.Context, userID, vehicleID int) error {
	query := `INSERT INTO favorites (user_id, vehicle_id) VALUES ($1, $2) ON CONFLICT (user_id, vehicle_id) DO NOTHING`
	_, err := r.db.ExecContext(ctx, query, userID, vehicleID)
	if err != nil {
		return fmt.Errorf("error adding to favorites: %w", err)
	}
	return nil
}

// RemoveFromFavorites удаляет технику из избранного
func (r *FavoriteRepository) RemoveFromFavorites(ctx context.Context, userID, vehicleID int) error {
	query := `DELETE FROM favorites WHERE user_id = $1 AND vehicle_id = $2`
	_, err := r.db.ExecContext(ctx, query, userID, vehicleID)
	if err != nil {
		return fmt.Errorf("error removing from favorites: %w", err)
	}
	return nil
}

// GetUserFavorites получает все избранные пользователя
func (r *FavoriteRepository) GetUserFavorites(ctx context.Context, userID int) ([]models.FavoriteWithVehicle, error) {
	query := `
		SELECT 
			f.favorite_id,
			f.user_id,
			f.vehicle_id,
			f.created_at,
			vm.model_name,
			vt.type_name as category,
			v.price,
			v.status,
			v.manufacture_year as year,
			v.color,
			v.serial_number,
			v.vin
		FROM favorites f
		JOIN vehicles v ON f.vehicle_id = v.vehicle_id
		JOIN vehicle_models vm ON v.model_id = vm.model_id
		JOIN vehicle_types vt ON vm.type_id = vt.type_id
		WHERE f.user_id = $1
		ORDER BY f.created_at DESC
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user favorites: %w", err)
	}
	defer rows.Close()
	
	var favorites []models.FavoriteWithVehicle
	
	for rows.Next() {
		var fav models.FavoriteWithVehicle
		err := rows.Scan(
			&fav.FavoriteID,
			&fav.UserID,
			&fav.VehicleID,
			&fav.CreatedAt,
			&fav.ModelName,
			&fav.Category,
			&fav.Price,
			&fav.Status,
			&fav.Year,
			&fav.Color,
			&fav.SerialNumber,
			&fav.VIN,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning favorite: %w", err)
		}
		favorites = append(favorites, fav)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating favorites: %w", err)
	}
	
	return favorites, nil
}

// IsFavorite проверяет, находится ли техника в избранном
func (r *FavoriteRepository) IsFavorite(ctx context.Context, userID, vehicleID int) (bool, error) {
	query := `SELECT COUNT(*) FROM favorites WHERE user_id = $1 AND vehicle_id = $2`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID, vehicleID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking if favorite: %w", err)
	}
	return count > 0, nil
}

// GetFavoriteCount получает количество избранных у пользователя
func (r *FavoriteRepository) GetFavoriteCount(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM favorites WHERE user_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("error getting favorite count: %w", err)
	}
	return count, nil
}
