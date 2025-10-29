package handlers

import (
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type FavoriteHandler struct {
	service *service.FavoriteService
}

func NewFavoriteHandler(service *service.FavoriteService) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

// GetUserFavorites получает все избранные пользователя
func (h *FavoriteHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из JWT токена
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	favorites, err := h.service.GetUserFavorites(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка получения избранного")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, favorites)
}

// ToggleFavorite переключает статус избранного
func (h *FavoriteHandler) ToggleFavorite(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из JWT токена
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	// Получаем vehicleID из URL
	vars := mux.Vars(r)
	vehicleIDStr, ok := vars["id"]
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "ID техники не указан")
		return
	}

	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Неверный ID техники")
		return
	}

	isFavorite, err := h.service.ToggleFavorite(r.Context(), userID, vehicleID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка изменения избранного")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"is_favorite": isFavorite,
		"message": func() string {
			if isFavorite {
				return "Техника добавлена в избранное"
			}
			return "Техника удалена из избранного"
		}(),
	})
}

// IsFavorite проверяет, находится ли техника в избранном
func (h *FavoriteHandler) IsFavorite(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из JWT токена
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	// Получаем vehicleID из URL
	vars := mux.Vars(r)
	vehicleIDStr, ok := vars["id"]
	if !ok {
		utils.ErrorResponse(w, http.StatusBadRequest, "ID техники не указан")
		return
	}

	vehicleID, err := strconv.Atoi(vehicleIDStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Неверный ID техники")
		return
	}

	isFavorite, err := h.service.IsFavorite(r.Context(), userID, vehicleID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка проверки избранного")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"is_favorite": isFavorite,
	})
}

// GetFavoriteCount получает количество избранных
func (h *FavoriteHandler) GetFavoriteCount(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из JWT токена
	userID, err := h.getUserIDFromToken(r)
	if err != nil {
		utils.ErrorResponse(w, http.StatusUnauthorized, "Недействительный токен")
		return
	}

	count, err := h.service.GetFavoriteCount(r.Context(), userID)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка получения количества избранных")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]interface{}{
		"count": count,
	})
}

// getUserIDFromToken извлекает userID из JWT токена
func (h *FavoriteHandler) getUserIDFromToken(r *http.Request) (int, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("токен авторизации не найден")
	}

	tokenString := authHeader
	if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
		tokenString = authHeader[7:]
	}

	claims, err := utils.ValidateJWTClaims(tokenString, "amkodor-secret-key-change-in-production")
	if err != nil {
		return 0, fmt.Errorf("недействительный токен: %w", err)
	}

	return claims.UserID, nil
}
