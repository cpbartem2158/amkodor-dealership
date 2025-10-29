package handlers

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"amkodor-dealership/internal/models"
	"amkodor-dealership/internal/service"
	"amkodor-dealership/internal/utils"

	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type ReportHandler struct {
	service *service.ReportService
}

func NewReportHandler(service *service.ReportService) *ReportHandler {
	return &ReportHandler{
		service: service,
	}
}

// SalesReport отчёт по продажам
func (h *ReportHandler) SalesReport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	startDate, _ := time.Parse("2006-01-02", query.Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", query.Get("end_date"))
	employeeID, _ := strconv.Atoi(query.Get("employee_id"))

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0) // последний месяц по умолчанию
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	_ = models.ReportFilters{
		StartDate:  startDate,
		EndDate:    endDate,
		EmployeeID: employeeID,
	}

	// Заглушка для GenerateSalesReport
	report := map[string]interface{}{
		"total_sales":   100,
		"total_revenue": 1000000,
		"period":        "2024-01-01 to 2024-01-31",
	}

	utils.SuccessResponse(w, http.StatusOK, report)
}

// ExportToExcel экспорт отчёта в Excel
func (h *ReportHandler) ExportToExcel(w http.ResponseWriter, r *http.Request) {
	var req models.ExportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Некорректные данные")
		return
	}

	// Заглушка для ExportToExcel
	excelFile := []byte("dummy excel content")
	filename := "report.xlsx"

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))

	if _, err := w.Write(excelFile); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка отправки файла")
		return
	}
}

// EmployeesReport отчёт по сотрудникам
func (h *ReportHandler) EmployeesReport(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	startDate, _ := time.Parse("2006-01-02", query.Get("start_date"))
	endDate, _ := time.Parse("2006-01-02", query.Get("end_date"))

	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	// Заглушка для GenerateEmployeesReport
	report := map[string]interface{}{
		"total_employees":  50,
		"active_employees": 45,
		"period":           "2024-01-01 to 2024-01-31",
	}

	utils.SuccessResponse(w, http.StatusOK, report)
}

// VehiclesReport отчёт по автомобилям
func (h *ReportHandler) VehiclesReport(w http.ResponseWriter, r *http.Request) {
	// Заглушка для GenerateVehiclesReport
	report := map[string]interface{}{
		"total_vehicles":     200,
		"available_vehicles": 150,
		"reserved_vehicles":  30,
		"sold_vehicles":      20,
	}

	utils.SuccessResponse(w, http.StatusOK, report)
}

// InventoryReport отчёт по инвентарю
func (h *ReportHandler) InventoryReport(w http.ResponseWriter, r *http.Request) {
	// Заглушка для InventoryReport
	report := map[string]interface{}{
		"total_items":     500,
		"available_items": 450,
		"reserved_items":  30,
		"low_stock_items": 20,
	}

	utils.SuccessResponse(w, http.StatusOK, report)
}

// ExportSalesReport экспорт отчёта по продажам
func (h *ReportHandler) ExportSalesReport(w http.ResponseWriter, r *http.Request) {
	// Заглушка для ExportSalesReport
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=sales_report.xlsx")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Sales report export not implemented"))
}

// ExportInventoryReport экспорт отчёта по инвентарю
func (h *ReportHandler) ExportInventoryReport(w http.ResponseWriter, r *http.Request) {
	// Заглушка для ExportInventoryReport
	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=inventory_report.xlsx")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Inventory report export not implemented"))
}

// ExportReport - единая точка экспорта: csv, xlsx, pdf
func (h *ReportHandler) ExportReport(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	reportType := q.Get("type") // sales, vehicles, customers, financial
	format := q.Get("format")   // csv | xlsx | pdf (может прийти excel/xslx)
	if format == "excel" || format == "xslx" {
		format = "xlsx"
	}
	startStr := q.Get("start_date")
	endStr := q.Get("end_date")

	startDate, _ := time.Parse("2006-01-02", startStr)
	endDate, _ := time.Parse("2006-01-02", endStr)
	if startDate.IsZero() {
		startDate = time.Now().AddDate(0, -1, 0)
	}
	if endDate.IsZero() {
		endDate = time.Now()
	}

	filenameBase := map[string]string{
		"sales":     "sales_report",
		"vehicles":  "vehicles_report",
		"customers": "customers_report",
		"financial": "financial_report",
	}[reportType]
	if filenameBase == "" {
		filenameBase = "report"
	}

	switch format {
	case "csv":
		buf := &bytes.Buffer{}
		wtr := csv.NewWriter(buf)
		// Примеры данных (можно заменить реальными)
		switch reportType {
		case "sales":
			wtr.Write([]string{"ID", "Клиент", "Техника", "Сумма", "Дата", "Статус"})
			wtr.Write([]string{"1", "Иван Петров", "Экскаватор", "450000", startDate.Format("2006-01-02"), "Завершена"})
		case "vehicles":
			wtr.Write([]string{"ID", "Название", "Категория", "Цена", "Статус", "Год"})
			wtr.Write([]string{"1", "Экскаватор 331.04", "Экскаваторы", "450000", "В наличии", "2023"})
		case "customers":
			wtr.Write([]string{"ID", "Название", "Тип", "Email", "Телефон"})
			wtr.Write([]string{"1", "Иван Петров", "ФЛ", "ivan@example.com", "+375291234567"})
		case "financial":
			wtr.Write([]string{"Период", "Доходы", "Расходы", "Прибыль", "Продажи"})
			wtr.Write([]string{startDate.Format("2006-01-02"), "450000", "50000", "400000", "1"})
		default:
			wtr.Write([]string{"Отчет", "Период"})
			wtr.Write([]string{"Общий", fmt.Sprintf("%s - %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))})
		}
		wtr.Flush()
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_%s.csv", filenameBase, startDate.Format("20060102"), endDate.Format("20060102")))
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return

	case "xlsx":
		f := excelize.NewFile()
		defer f.Close()
		sheet := "Отчет"
		idx, _ := f.NewSheet(sheet)
		f.SetActiveSheet(idx)
		f.SetCellValue(sheet, "A1", fmt.Sprintf("Отчет: %s", reportType))
		f.SetCellValue(sheet, "A2", fmt.Sprintf("Период: %s - %s", startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))
		// Заголовки и строки по типу
		switch reportType {
		case "sales":
			headers := []string{"ID", "Клиент", "Техника", "Сумма", "Дата", "Статус"}
			for i, hname := range headers {
				f.SetCellValue(sheet, fmt.Sprintf("%c4", 'A'+i), hname)
			}
			f.SetCellValue(sheet, "A5", 1)
			f.SetCellValue(sheet, "B5", "Иван Петров")
			f.SetCellValue(sheet, "C5", "Экскаватор")
			f.SetCellValue(sheet, "D5", 450000)
			f.SetCellValue(sheet, "E5", startDate.Format("2006-01-02"))
			f.SetCellValue(sheet, "F5", "Завершена")
		default:
			f.SetCellValue(sheet, "A4", "Данные")
			f.SetCellValue(sheet, "A5", "Пример строки")
		}
		buf, err := f.WriteToBuffer()
		if err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка формирования XLSX")
			return
		}
		w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_%s.xlsx", filenameBase, startDate.Format("20060102"), endDate.Format("20060102")))
		w.WriteHeader(http.StatusOK)
		w.Write(buf.Bytes())
		return

	case "pdf":
		pdf := gofpdf.New("P", "mm", "A4", "")
		// Подключаем UTF-8 шрифты для корректной кириллицы (из файлов)
		fontReg, _ := filepath.Abs("web/static/fonts/DejaVuSans.ttf")
		fontBold, _ := filepath.Abs("web/static/fonts/DejaVuSans-Bold.ttf")
		pdf.AddUTF8Font("DejaVu", "", fontReg)
		pdf.AddUTF8Font("DejaVu", "B", fontBold)
		if err := pdf.Error(); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка загрузки шрифта для PDF")
			return
		}
		pdf.AddPage()
		pdf.SetFont("DejaVu", "B", 16)
		pdf.Cell(190, 10, fmt.Sprintf("Отчет: %s", reportType))
		pdf.Ln(8)
		pdf.SetFont("DejaVu", "", 12)
		pdf.Cell(190, 8, fmt.Sprintf("Период: %s - %s", startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))
		pdf.Ln(12)
		// Простая таблица-заглушка
		pdf.CellFormat(30, 8, "Колонка 1", "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, "Колонка 2", "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, "Колонка 3", "1", 1, "L", false, 0, "")
		pdf.CellFormat(30, 8, "Значение 1", "1", 0, "L", false, 0, "")
		pdf.CellFormat(60, 8, "Значение 2", "1", 0, "L", false, 0, "")
		pdf.CellFormat(30, 8, "Значение 3", "1", 1, "L", false, 0, "")
		var out bytes.Buffer
		if err := pdf.Output(&out); err != nil {
			utils.ErrorResponse(w, http.StatusInternalServerError, "Ошибка формирования PDF")
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s_%s_%s.pdf", filenameBase, startDate.Format("20060102"), endDate.Format("20060102")))
		w.WriteHeader(http.StatusOK)
		w.Write(out.Bytes())
		return
	default:
		utils.ErrorResponse(w, http.StatusBadRequest, "Неверный формат отчёта")
		return
	}
}
