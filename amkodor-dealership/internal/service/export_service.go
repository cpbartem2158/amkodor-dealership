package service

import (
	"amkodor-dealership/internal/models"
	"database/sql"
	"fmt"
)

type ExportService struct {
	db *sql.DB
}

func NewExportService(db *sql.DB) *ExportService {
	return &ExportService{db: db}
}

func (s *ExportService) ExportSalesToExcel(sales []models.Sale) (string, error) {
	// Реализация экспорта будет в отдельном пакете excel
	return "", fmt.Errorf("not implemented yet")
}

func (s *ExportService) ExportInventoryToExcel(inventory []map[string]interface{}) (string, error) {
	// Реализация экспорта будет в отдельном пакете excel
	return "", fmt.Errorf("not implemented yet")
}

func (s *ExportService) LogExport(packageName string, rowsProcessed int, status, errorMessage string) error {
	query := `
		INSERT INTO ssis_log (package_name, rows_processed, status, error_message, username, execution_duration)
		VALUES ($1, $2, $3, $4, CURRENT_USER, 0)
	`

	_, err := s.db.Exec(query, packageName, rowsProcessed, status, errorMessage)
	if err != nil {
		return fmt.Errorf("error logging export: %w", err)
	}

	return nil
}
