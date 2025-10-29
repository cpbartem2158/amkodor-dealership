package repository

import (
	"context"
	"database/sql"
	"fmt"
	"amkodor-dealership/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(ctx context.Context, user *models.User) (int, error) {
	query := `INSERT INTO users (name, email, phone, password_hash, role) 
			  VALUES ($1, $2, $3, $4, $5) RETURNING user_id`
	
	fmt.Printf("Creating user: %s, %s, %s\n", user.Name, user.Email, user.Phone)
	
	var userID int
	err := r.db.QueryRow(query, 
		user.Name, user.Email, user.Phone, user.PasswordHash, user.Role).Scan(&userID)
	
	if err != nil {
		fmt.Printf("Error creating user: %v\n", err)
		return 0, fmt.Errorf("error creating user: %w", err)
	}
	
	fmt.Printf("User created with ID: %d\n", userID)
	return userID, nil
}

// GetByEmail получает пользователя по email
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT user_id, name, email, phone, password_hash, role, created_at, updated_at 
			  FROM users WHERE email = $1`
	
	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.UserID, &user.Name, &user.Email, &user.Phone, 
		&user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	
	return &user, nil
}

// GetByID получает пользователя по ID
func (r *UserRepository) GetByID(ctx context.Context, userID int) (*models.User, error) {
	query := `SELECT user_id, name, email, phone, password_hash, role, created_at, updated_at 
			  FROM users WHERE user_id = $1`
	
	var user models.User
	err := r.db.QueryRow(query, userID).Scan(
		&user.UserID, &user.Name, &user.Email, &user.Phone, 
		&user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("error getting user: %w", err)
	}
	
	return &user, nil
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(ctx context.Context, user *models.User) error {
	query := `UPDATE users SET name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP 
			  WHERE user_id = $3`
	
	_, err := r.db.Exec(query, user.Name, user.Phone, user.UserID)
	if err != nil {
		return fmt.Errorf("error updating user: %w", err)
	}
	
	return nil
}

// Delete удаляет пользователя
func (r *UserRepository) Delete(ctx context.Context, userID int) error {
	query := `DELETE FROM users WHERE user_id = $1`
	
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("error deleting user: %w", err)
	}
	
	return nil
}

// EmailExists проверяет, существует ли пользователь с таким email
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = $1`
	
	var count int
	err := r.db.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking email: %w", err)
	}
	
	return count > 0, nil
}
