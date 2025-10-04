package excel

import (
	"fmt"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExportSalesReport экспортирует отчет о продажах в Excel
func ExportSalesReport(sales interface{}, startDate, endDate time.Time) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Продажи"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	// Установка активного листа
	f.SetActiveSheet(index)

	// Заголовок отчета
	f.SetCellValue(sheetName, "A1", "Отчет по продажам")
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Период: %s - %s", startDate.Format("02.01.2006"), endDate.Format("02.01.2006")))

	// Заголовки таблицы
	headers := []string{"№", "Дата", "Договор", "Модель", "Клиент", "Менеджер", "Базовая цена", "Скидка", "Итого", "Оплата", "Статус"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c4", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Стиль заголовков
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#4472C4"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle(sheetName, "A4", "K4", headerStyle)

	// TODO: Добавить данные продаж
	// Пример:
	// for i, sale := range sales {
	//     row := i + 5
	//     f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), i+1)
	//     f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), sale.Date)
	//     ...
	// }

	// Автоширина колонок
	f.SetColWidth(sheetName, "A", "K", 15)

	// Сохранение файла
	filename := fmt.Sprintf("sales_report_%s.xlsx", time.Now().Format("20060102_150405"))
	filepath := fmt.Sprintf("./exports/%s", filename)

	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return filepath, nil
}

// ExportInventoryReport экспортирует отчет по инвентарю в Excel
func ExportInventoryReport(inventory interface{}) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Инвентарь"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	f.SetActiveSheet(index)

	// Заголовок отчета
	f.SetCellValue(sheetName, "A1", "Отчет по инвентарю")
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Дата: %s", time.Now().Format("02.01.2006 15:04")))

	// Заголовки таблицы
	headers := []string{"Филиал", "Город", "Категория", "Тип", "Модель", "Количество", "Статус", "Средняя цена", "Общая стоимость"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c4", 'A'+i)
		f.SetCellValue(sheetName, cell, header)
	}

	// Стиль заголовков
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 12,
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#70AD47"},
			Pattern: 1,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	f.SetCellStyle(sheetName, "A4", "I4", headerStyle)

	// Автоширина колонок
	f.SetColWidth(sheetName, "A", "I", 15)

	// Сохранение файла
	filename := fmt.Sprintf("inventory_report_%s.xlsx", time.Now().Format("20060102_150405"))
	filepath := fmt.Sprintf("./exports/%s", filename)

	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return filepath, nil
}

// CreateExcelReport создает общий Excel отчет
func CreateExcelReport(data map[string]interface{}, reportType string) (string, error) {
	f := excelize.NewFile()
	defer f.Close()

	sheetName := "Отчет"
	index, err := f.NewSheet(sheetName)
	if err != nil {
		return "", fmt.Errorf("error creating sheet: %w", err)
	}

	f.SetActiveSheet(index)

	// Заголовок
	f.SetCellValue(sheetName, "A1", fmt.Sprintf("Отчет: %s", reportType))
	f.SetCellValue(sheetName, "A2", fmt.Sprintf("Дата формирования: %s", time.Now().Format("02.01.2006 15:04:05")))

	// Стили
	titleStyle, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{
			Bold: true,
			Size: 14,
		},
	})
	f.SetCellStyle(sheetName, "A1", "A1", titleStyle)

	// Сохранение
	filename := fmt.Sprintf("%s_%s.xlsx", reportType, time.Now().Format("20060102_150405"))
	filepath := fmt.Sprintf("./exports/%s", filename)

	if err := f.SaveAs(filepath); err != nil {
		return "", fmt.Errorf("error saving file: %w", err)
	}

	return filepath, nil
}
