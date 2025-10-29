package handlers

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"amkodor-dealership/internal/utils"
)

// Глобальное хранилище для сервисных заявок (временное решение)
var (
	serviceRequestsStore = make(map[int]map[string]interface{})
	serviceRequestsMutex sync.RWMutex
	nextServiceRequestID = 1
)

type UserHandler struct {
	// Здесь можно добавить зависимости от сервисов
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// GetUserStats получение статистики пользователя
func (h *UserHandler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	// Заглушка для демонстрации
	stats := map[string]interface{}{
		"totalOrders":   3,
		"favoriteItems": 5,
		"serviceVisits": 2,
		"totalSpent":    450000,
	}

	utils.RespondSuccess(w, stats)
}

// GetUserOrders получение заказов пользователя
func (h *UserHandler) GetUserOrders(w http.ResponseWriter, r *http.Request) {
	// Заглушка для демонстрации
	orders := []map[string]interface{}{
		{
			"id":     1,
			"number": "ORD-2024-001",
			"date":   "2024-01-15",
			"status": "processing",
			"items": []map[string]interface{}{
				{
					"name":  "Экскаватор АМКОДОР 3316",
					"specs": "2023 год, Дизель, 160 л.с.",
					"price": 150000,
					"image": "🚜",
				},
			},
			"total": 150000,
		},
		{
			"id":     2,
			"number": "ORD-2024-002",
			"date":   "2024-01-10",
			"status": "delivered",
			"items": []map[string]interface{}{
				{
					"name":  "Бульдозер АМКОДОР 3317",
					"specs": "2023 год, Дизель, 180 л.с.",
					"price": 200000,
					"image": "🚜",
				},
			},
			"total": 200000,
		},
	}

	utils.RespondSuccess(w, orders)
}

// GetUserFavorites получение избранного пользователя
func (h *UserHandler) GetUserFavorites(w http.ResponseWriter, r *http.Request) {
	// Заглушка для демонстрации
	favorites := []map[string]interface{}{
		{
			"id":       1,
			"name":     "Экскаватор АМКОДОР 3316",
			"category": "Экскаваторы",
			"price":    150000,
			"year":     2025,
			"engine":   "Дизель",
			"power":    "160 л.с.",
			"weight":   "16 т",
			"image":    "🚜",
		},
		{
			"id":       2,
			"name":     "Бульдозер АМКОДОР 3317",
			"category": "Бульдозеры",
			"price":    200000,
			"year":     2025,
			"engine":   "Дизель",
			"power":    "180 л.с.",
			"weight":   "17 т",
			"image":    "🚜",
		},
	}

	utils.RespondSuccess(w, favorites)
}

// AddToFavorites добавление в избранное
func (h *UserHandler) AddToFavorites(w http.ResponseWriter, r *http.Request) {
	var req struct {
		VehicleID int `json:"vehicle_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	// Заглушка для демонстрации
	utils.RespondSuccess(w, map[string]string{
		"message": "Техника добавлена в избранное",
	})
}

// RemoveFromFavorites удаление из избранного
func (h *UserHandler) RemoveFromFavorites(w http.ResponseWriter, r *http.Request) {
	vehicleID := r.URL.Query().Get("vehicle_id")
	if vehicleID == "" {
		utils.RespondError(w, http.StatusBadRequest, "ID техники обязателен")
		return
	}

	// Заглушка для демонстрации
	utils.RespondSuccess(w, map[string]string{
		"message": "Техника удалена из избранного",
	})
}

// GetServiceRequests получение заявок на сервис
func (h *UserHandler) GetServiceRequests(w http.ResponseWriter, r *http.Request) {
	// Базовые тестовые данные
	baseRequests := []map[string]interface{}{
		{
			"id":          1,
			"title":       "Техническое обслуживание экскаватора",
			"description": "Плановое ТО экскаватора АМКОДОР 3316",
			"date":        "2025-01-15",
			"status":      "pending",
		},
		{
			"id":          2,
			"title":       "Ремонт гидравлики",
			"description": "Замена гидроцилиндра подъема стрелы",
			"date":        "2025-01-10",
			"status":      "in-progress",
		},
	}

	// Получаем созданные пользователями заявки
	serviceRequestsMutex.RLock()
	userRequests := make([]map[string]interface{}, 0, len(serviceRequestsStore))
	for _, request := range serviceRequestsStore {
		userRequests = append(userRequests, request)
	}
	serviceRequestsMutex.RUnlock()

	// Объединяем базовые и пользовательские заявки
	allRequests := append(baseRequests, userRequests...)

	utils.RespondSuccess(w, allRequests)
}

// CreateServiceRequest создание заявки на сервис
func (h *UserHandler) CreateServiceRequest(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		ServiceType string `json:"service_type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	// Сохраняем заявку в глобальное хранилище
	serviceRequestsMutex.Lock()
	id := nextServiceRequestID
	nextServiceRequestID++

	serviceRequestsStore[id] = map[string]interface{}{
		"id":           id,
		"title":        req.Title,
		"description":  req.Description,
		"service_type": req.ServiceType,
		"date":         time.Now().Format("2006-01-02"),
		"status":       "pending",
	}
	serviceRequestsMutex.Unlock()

	utils.RespondSuccess(w, map[string]interface{}{
		"id":      id,
		"message": "Заявка на сервис создана",
	})
}

// GetAllServiceRequests получение всех сервисных заявок для админа
func (h *UserHandler) GetAllServiceRequests(w http.ResponseWriter, r *http.Request) {
	// Базовые тестовые данные
	baseRequests := []map[string]interface{}{
		{
			"id":          1,
			"title":       "Техническое обслуживание экскаватора",
			"description": "Плановое ТО экскаватора АМКОДОР 3316",
			"date":        "2025-01-15",
			"status":      "pending",
			"customer":    "Иван Иванов",
			"phone":       "+375 29 123-45-67",
		},
		{
			"id":          2,
			"title":       "Ремонт гидравлики",
			"description": "Замена гидроцилиндра подъема стрелы",
			"date":        "2025-01-10",
			"status":      "in-progress",
			"customer":    "Петр Петров",
			"phone":       "+375 29 234-56-78",
		},
	}

	// Получаем созданные пользователями заявки
	serviceRequestsMutex.RLock()
	userRequests := make([]map[string]interface{}, 0, len(serviceRequestsStore))
	for _, request := range serviceRequestsStore {
		// Добавляем информацию о клиенте для админа
		adminRequest := make(map[string]interface{})
		for k, v := range request {
			adminRequest[k] = v
		}
		adminRequest["customer"] = "Клиент"
		adminRequest["phone"] = "Не указан"
		userRequests = append(userRequests, adminRequest)
	}
	serviceRequestsMutex.RUnlock()

	// Объединяем базовые и пользовательские заявки
	allRequests := append(baseRequests, userRequests...)

	utils.RespondSuccess(w, allRequests)
}

// GetUserProfile получение профиля пользователя
func (h *UserHandler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	// Заглушка для демонстрации
	profile := map[string]interface{}{
		"id":        1,
		"name":      "Иван Иванов",
		"email":     "ivan@example.com",
		"phone":     "+375 29 123-45-67",
		"company":   "ООО Стройтехника",
		"createdAt": "2024-01-01",
	}

	utils.RespondSuccess(w, profile)
}

// UpdateUserProfile обновление профиля пользователя
func (h *UserHandler) UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	// Заглушка для демонстрации
	utils.RespondSuccess(w, map[string]string{
		"message": "Профиль обновлен",
	})
}
