package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"
)

// Глобальное хранилище для техники (временное решение)
var (
	vehiclesStore = make(map[int]models.CreateVehicleRequest)
	vehiclesMutex sync.RWMutex
	nextVehicleID = 1
)

type VehicleHandler struct {
	service *service.VehicleService
}

func NewVehicleHandler(service *service.VehicleService) *VehicleHandler {
	return &VehicleHandler{
		service: service,
	}
}

// GetAll получение всех автомобилей
func (h *VehicleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	log.Println("GetAll: получен запрос на получение всех автомобилей")

	// Базовые тестовые данные
	baseVehicles := []map[string]interface{}{
		{
			"id":       1,
			"name":     "Экскаватор АМКОДОР 3316",
			"category": "Экскаваторы",
			"price":    150000,
			"status":   "В наличии",
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
			"status":   "В наличии",
			"year":     2025,
			"engine":   "Дизель",
			"power":    "180 л.с.",
			"weight":   "17 т",
			"image":    "🚜",
		},
		{
			"id":       3,
			"name":     "Погрузчик АМКОДОР 3318",
			"category": "Погрузчики",
			"price":    120000,
			"status":   "В наличии",
			"year":     2025,
			"engine":   "Дизель",
			"power":    "140 л.с.",
			"weight":   "14 т",
			"image":    "🚜",
		},
		{
			"id":       4,
			"name":     "Кран АМКОДОР 3319",
			"category": "Краны",
			"price":    300000,
			"status":   "В наличии",
			"year":     2025,
			"engine":   "Дизель",
			"power":    "200 л.с.",
			"weight":   "25 т",
			"image":    "🚜",
		},
	}

	// Получаем добавленную админом технику
	vehiclesMutex.RLock()
	addedVehicles := make([]map[string]interface{}, 0, len(vehiclesStore))
	for id, vehicle := range vehiclesStore {
		addedVehicles = append(addedVehicles, map[string]interface{}{
			"id":       id,
			"name":     vehicle.Name,
			"category": vehicle.Category,
			"price":    vehicle.Price,
			"status":   vehicle.Status,
			"year":     vehicle.Year,
			"engine":   vehicle.Engine,
			"power":    vehicle.Power,
			"weight":   vehicle.Weight,
			"image":    vehicle.Image,
		})
	}
	vehiclesMutex.RUnlock()

	// Объединяем базовые и добавленные данные
	allVehicles := append(baseVehicles, addedVehicles...)

	log.Printf("GetAll: возвращаем %d автомобилей (%d базовых + %d добавленных)",
		len(allVehicles), len(baseVehicles), len(addedVehicles))

	utils.RespondSuccess(w, allVehicles)
}

// GetByID получение автомобиля по ID
func (h *VehicleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	vehicle, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Автомобиль не найден")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, vehicle)
}

// Create создание нового автомобиля
func (h *VehicleHandler) Create(w http.ResponseWriter, r *http.Request) {
	log.Println("Create: получен запрос на создание техники")

	var req models.CreateVehicleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	log.Printf("Create: получены данные: %+v", req)

	// Сохраняем в глобальное хранилище
	vehiclesMutex.Lock()
	id := nextVehicleID
	nextVehicleID++
	vehiclesStore[id] = req
	vehiclesMutex.Unlock()

	log.Printf("Create: техника сохранена с ID %d", id)

	utils.RespondSuccess(w, map[string]interface{}{
		"id":      id,
		"message": "Техника успешно создана",
		"data":    req,
	})
}

// Update обновление автомобиля
func (h *VehicleHandler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	var vehicle models.Vehicle
	if err := json.NewDecoder(r.Body).Decode(&vehicle); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	vehicle.VehicleID = id

	if err := utils.ValidateStruct(&vehicle); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Ошибка валидации")
		return
	}

	if err := h.service.Update(&vehicle); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка обновления автомобиля")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Автомобиль успешно обновлён",
	})
}

// Delete удаление автомобиля
func (h *VehicleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректный ID")
		return
	}

	if err := h.service.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка удаления автомобиля")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, map[string]string{
		"message": "Автомобиль успешно удалён",
	})
}

// Search поиск автомобилей по критериям
func (h *VehicleHandler) Search(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	filters := models.VehicleFilters{
		Model:     query.Get("model"),
		YearFrom:  parseIntParam(query.Get("year_from")),
		YearTo:    parseIntParam(query.Get("year_to")),
		PriceFrom: parseFloatParam(query.Get("price_from")),
		PriceTo:   parseFloatParam(query.Get("price_to")),
		Status:    query.Get("status"),
		Color:     query.Get("color"),
	}

	vehicles, err := h.service.Search(map[string]interface{}{
		"model":      filters.Model,
		"year_from":  filters.YearFrom,
		"year_to":    filters.YearTo,
		"price_from": filters.PriceFrom,
		"price_to":   filters.PriceTo,
		"status":     filters.Status,
		"color":      filters.Color,
	})
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка поиска автомобилей")
		return
	}

	utils.SuccessResponse(w, http.StatusOK, vehicles)
}

// GetCategories получение всех категорий
func (h *VehicleHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	// Заглушка для GetCategories
	categories := []string{"Экскаваторы", "Бульдозеры", "Погрузчики", "Тракторы"}

	utils.SuccessResponse(w, http.StatusOK, categories)
}

// Helper функции
func parseIntParam(param string) *int {
	if param == "" {
		return nil
	}
	val, err := strconv.Atoi(param)
	if err != nil {
		return nil
	}
	return &val
}

func parseFloatParam(param string) *float64 {
	if param == "" {
		return nil
	}
	val, err := strconv.ParseFloat(param, 64)
	if err != nil {
		return nil
	}
	return &val
}

// UploadImage обработка загрузки изображения
func (h *VehicleHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	log.Println("UploadImage: получен запрос на загрузку изображения")

	// Добавляем CORS заголовки
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Обрабатываем preflight запросы
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Проверяем метод
	if r.Method != http.MethodPost {
		utils.RespondError(w, http.StatusMethodNotAllowed, "Метод не поддерживается")
		return
	}

	// Парсим multipart form
	err := r.ParseMultipartForm(10 << 20) // 10 MB
	if err != nil {
		log.Printf("UploadImage: ошибка парсинга формы: %v", err)
		utils.RespondError(w, http.StatusBadRequest, "Ошибка парсинга формы")
		return
	}

	// Получаем файл
	file, handler, err := r.FormFile("image")
	if err != nil {
		log.Printf("UploadImage: файл не найден: %v", err)
		utils.RespondError(w, http.StatusBadRequest, "Файл не найден")
		return
	}
	defer file.Close()

	log.Printf("UploadImage: получен файл %s, размер %d байт", handler.Filename, handler.Size)

	// Проверяем тип файла
	contentType := handler.Header.Get("Content-Type")
	if len(contentType) < 5 || contentType[:5] != "image" {
		utils.RespondError(w, http.StatusBadRequest, "Файл должен быть изображением")
		return
	}

	// Создаем уникальное имя файла
	ext := filepath.Ext(handler.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	// Создаем папку для загрузок, если её нет
	uploadDir := "web/static/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка создания папки: "+err.Error())
		return
	}

	// Создаем файл на диске
	filepath := filepath.Join(uploadDir, filename)
	dst, err := os.Create(filepath)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка создания файла")
		return
	}
	defer dst.Close()

	// Копируем содержимое файла
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, "Ошибка сохранения файла")
		return
	}

	// Возвращаем URL изображения
	imageURL := fmt.Sprintf("/static/uploads/%s", filename)
	utils.RespondSuccess(w, map[string]string{
		"url":      imageURL,
		"filename": filename,
	})
}
