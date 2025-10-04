package service

import (
	"amkodor-dealership/internal/models"
	"database/sql"
	"fmt"
	"time"
)

type ReportService struct {
	db *sql.DB
}

func NewReportService(db *sql.DB) *ReportService {
	return &ReportService{db: db}
}

func (s *ReportService) GenerateSalesReport(startDate, endDate time.Time, employeeID, warehouseID, categoryID *int) ([]models.Sale, error) {
	query := `SELECT * FROM sp_generate_sales_report($1, $2, $3, $4, $5)`

	rows, err := s.db.Query(query, startDate, endDate, employeeID, warehouseID, categoryID)
	if err != nil {
		return nil, fmt.Errorf("error generating sales report: %w", err)
	}
	defer rows.Close()

	sales := []models.Sale{}
	for rows.Next() {
		var sale models.Sale
		err := rows.Scan(
			&sale.SaleID, &sale.ContractNumber, &sale.SaleDate, &sale.ModelName,
			&sale.TypeName, &sale.CategoryName, &sale.ClientName, &sale.ManagerName,
			&sale.WarehouseName, &sale.WarehouseCity, &sale.BasePrice,
			&sale.DiscountAmount, &sale.FinalPrice, &sale.PaymentType, &sale.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning sale: %w", err)
		}
		sales = append(sales, sale)
	}

	return sales, nil
}

func (s *ReportService) GenerateInventoryReport(warehouseID, categoryID *int, status *string) ([]map[string]interface{}, error) {
	query := `SELECT * FROM sp_generate_inventory_report($1, $2, $3)`

	rows, err := s.db.Query(query, warehouseID, categoryID, status)
	if err != nil {
		return nil, fmt.Errorf("error generating inventory report: %w", err)
	}
	defer rows.Close()

	report := []map[string]interface{}{}
	for rows.Next() {
		var warehouseName, city, categoryName, typeName, modelName, vehicleStatus string
		var quantity int64
		var avgPrice, totalValue float64
		var oldestDate, newestDate time.Time

		err := rows.Scan(
			&warehouseName, &city, &categoryName, &typeName, &modelName,
			&quantity, &vehicleStatus, &avgPrice, &totalValue, &oldestDate, &newestDate,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning inventory: %w", err)
		}

		report = append(report, map[string]interface{}{
			"warehouse_name": warehouseName,
			"city":           city,
			"category_name":  categoryName,
			"type_name":      typeName,
			"model_name":     modelName,
			"quantity":       quantity,
			"status":         vehicleStatus,
			"average_price":  avgPrice,
			"total_value":    totalValue,
			"oldest_date":    oldestDate,
			"newest_date":    newestDate,
		})
	}

	return report, nil
}

func (s *ReportService) GetSalesStatistics(startDate, endDate time.Time) (map[string]interface{}, error) {
	query := `SELECT * FROM sp_get_sales_statistics($1, $2)`

	var totalSales int64
	var totalRevenue, avgSalePrice, totalDiscounts float64
	var bestCategory, bestManager string

	err := s.db.QueryRow(query, startDate, endDate).Scan(
		&totalSales, &totalRevenue, &avgSalePrice, &totalDiscounts, &bestCategory, &bestManager,
	)

	if err != nil {
		return nil, fmt.Errorf("error getting sales statistics: %w", err)
	}

	return map[string]interface{}{
		"total_sales":        totalSales,
		"total_revenue":      totalRevenue,
		"average_sale_price": avgSalePrice,
		"total_discounts":    totalDiscounts,
		"best_category":      bestCategory,
		"best_manager":       bestManager,
	}, nil
}
