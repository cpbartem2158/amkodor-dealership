package handlers

import (
	"encoding/json"
	"net/http"

	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{
		service: service,
	}
}

// Login авторизация пользователя
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if req.Email == "" || req.Password == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Email и пароль обязательны")
		return
	}

	// Реальная авторизация через сервис
	user, token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Неверные учетные данные")
		return
	}

	utils.RespondSuccess(w, map[string]interface{}{
		"token":   token,
		"message": "Login successful",
		"name":    user.Name,
		"role":    user.Role,
	})
}

// Register регистрация нового пользователя
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	if err := utils.ValidateStruct(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Ошибка валидации")
		return
	}

	// Заглушка для Register
	id := 1

	utils.SuccessResponse(w, http.StatusCreated, map[string]interface{}{
		"id":      id,
		"message": "Пользователь успешно зарегистрирован",
	})
}

// Logout выход из системы
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// В случае JWT токенов, logout обычно обрабатывается на клиенте
	// удалением токена из localStorage
	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Успешный выход",
	})
}

// GetCurrentUser получение информации о текущем пользователе
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value("userID").(int)
	if !ok {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Пользователь не авторизован")
		return
	}

	// Заглушка для GetUserByID
	user := map[string]interface{}{
		"id":    userID,
		"email": "user@example.com",
		"name":  "Test User",
	}

	utils.SuccessResponse(w, http.StatusOK, user)
}