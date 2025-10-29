// Модуль работы с отчётами

const Reports = {
    // Отчёт по продажам
    async salesReport(filters) {
        try {
            const queryString = new URLSearchParams(filters).toString();
            const response = await API.get(`/reports/sales?${queryString}`);
            return response;
        } catch (error) {
            console.error('Error getting sales report:', error);
            throw error;
        }
    },

    // Отчёт по сотрудникам
    async employeesReport(filters) {
        try {
            const queryString = new URLSearchParams(filters).toString();
            const response = await API.get(`/reports/employees?${queryString}`);
            return response;
        } catch (error) {
            console.error('Error getting employees report:', error);
            throw error;
        }
    },

    // Отчёт по автомобилям
    async vehiclesReport() {
        try {
            const response = await API.get('/reports/vehicles');
            return response;
        } catch (error) {
            console.error('Error getting vehicles report:', error);
            throw error;
        }
    },

    // Экспорт в Excel
    async exportToExcel(reportType, filters) {
        try {
            const response = await fetch('/api/reports/export/excel', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${Auth.getToken()}`
                },
                body: JSON.stringify({
                    report_type: reportType,
                    filters: filters
                })
            });

            if (response.ok) {
                const blob = await response.blob();
                const url = window.URL.createObjectURL(blob);
                const a = document.createElement('a');
                a.href = url;
                a.download = `report_${reportType}_${new Date().toISOString().split('T')[0]}.xlsx`;
                document.body.appendChild(a);
                a.click();
                a.remove();
                window.URL.revokeObjectURL(url);
                return { success: true };
            } else {
                return { success: false, error: 'Ошибка экспорта' };
            }
        } catch (error) {
            console.error('Error exporting to Excel:', error);
            return { success: false, error: 'Ошибка экспорта' };
        }
    }
};

// UI функции для работы с отчётами
const ReportsUI = {
    // Отобразить отчёт по продажам
    displaySalesReport(report, containerId = 'reportContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!report || !report.rows || report.rows.length === 0) {
            container.innerHTML = '<div class="no-results">Данные не найдены</div>';
            return;
        }

        const totalRevenue = report.rows.reduce((sum, row) => sum + row.final_price, 0);
        const totalProfit = report.rows.reduce((sum, row) => sum + row.profit, 0);
        const totalDiscount = report.rows.reduce((sum, row) => sum + row.discount, 0);

        const tableHTML = `
            <div class="report-summary">
                <div class="summary-card">
                    <h3>Общая выручка</h3>
                    <p class="summary-value">${window.app.formatCurrency(totalRevenue)}</p>
                </div>
                <div class="summary-card">
                    <h3>Общая прибыль</h3>
                    <p class="summary-value">${window.app.formatCurrency(totalProfit)}</p>
                </div>
                <div class="summary-card">
                    <h3>Общие скидки</h3>
                    <p class="summary-value">${window.app.formatCurrency(totalDiscount)}</p>
                </div>
                <div class="summary-card">
                    <h3>Количество продаж</h3>
                    <p class="summary-value">${report.rows.length}</p>
                </div>
            </div>

            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Дата</th>
                        <th>Автомобиль</th>
                        <th>VIN</th>
                        <th>Клиент</th>
                        <th>Сотрудник</th>
                        <th>Цена авто</th>
                        <th>Скидка</th>
                        <th>Итого</th>
                        <th>Прибыль</th>
                    </tr>
                </thead>
                <tbody>
                    ${report.rows.map(row => `
                        <tr>
                            <td>${row.id}</td>
                            <td>${window.app.formatDate(row.sale_date)}</td>
                            <td>${row.vehicle_model} (${row.year})</td>
                            <td>${row.vin}</td>
                            <td>${row.customer_name}</td>
                            <td>${row.employee_name}</td>
                            <td>${window.app.formatCurrency(row.vehicle_price)}</td>
                            <td>${window.app.formatCurrency(row.discount)}</td>
                            <td>${window.app.formatCurrency(row.final_price)}</td>
                            <td class="${row.profit >= 0 ? 'text-success' : 'text-danger'}">
                                ${window.app.formatCurrency(row.profit)}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Отобразить отчёт по сотрудникам
    displayEmployeesReport(report, containerId = 'reportContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!report || !report.rows || report.rows.length === 0) {
            container.innerHTML = '<div class="no-results">Данные не найдены</div>';
            return;
        }

        const totalRevenue = report.rows.reduce((sum, row) => sum + row.total_revenue, 0);
        const totalCommission = report.rows.reduce((sum, row) => sum + row.total_commission, 0);

        const tableHTML = `
            <div class="report-summary">
                <div class="summary-card">
                    <h3>Общая выручка</h3>
                    <p class="summary-value">${window.app.formatCurrency(totalRevenue)}</p>
                </div>
                <div class="summary-card">
                    <h3>Общие комиссионные</h3>
                    <p class="summary-value">${window.app.formatCurrency(totalCommission)}</p>
                </div>
            </div>

            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>ФИО</th>
                        <th>Должность</th>
                        <th>Количество продаж</th>
                        <th>Выручка</th>
                        <th>Ставка комиссии</th>
                        <th>Комиссионные</th>
                    </tr>
                </thead>
                <tbody>
                    ${report.rows.map(row => `
                        <tr>
                            <td>${row.id}</td>
                            <td>${row.full_name}</td>
                            <td>${row.position}</td>
                            <td>${row.sales_count}</td>
                            <td>${window.app.formatCurrency(row.total_revenue)}</td>
                            <td>${row.commission_rate}%</td>
                            <td>${window.app.formatCurrency(row.total_commission)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Отобразить отчёт по автомобилям
    displayVehiclesReport(report, containerId = 'reportContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!report || !report.rows || report.rows.length === 0) {
            container.innerHTML = '<div class="no-results">Данные не найдены</div>';
            return;
        }

        const tableHTML = `
            <table class="table">
                <thead>
                    <tr>
                        <th>Категория</th>
                        <th>Всего автомобилей</th>
                        <th>В наличии</th>
                        <th>Продано</th>
                        <th>Средняя цена</th>
                    </tr>
                </thead>
                <tbody>
                    ${report.rows.map(row => `
                        <tr>
                            <td>${row.category}</td>
                            <td>${row.total_count}</td>
                            <td>${row.available_count}</td>
                            <td>${row.sold_count}</td>
                            <td>${window.app.formatCurrency(row.avg_price)}</td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Генерация отчёта
    async generateReport(reportType) {
        try {
            window.app.showLoading();

            const filters = this.getFilters(reportType);
            let report;

            switch (reportType) {
                case 'sales':
                    report = await Reports.salesReport(filters);
                    this.displaySalesReport(report);
                    break;
                case 'employees':
                    report = await Reports.employeesReport(filters);
                    this.displayEmployeesReport(report);
                    break;
                case 'vehicles':
                    report = await Reports.vehiclesReport();
                    this.displayVehiclesReport(report);
                    break;
                default:
                    window.app.showNotification('Неизвестный тип отчёта', 'error');
            }
        } catch (error) {
            window.app.showNotification('Ошибка генерации отчёта', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Получить фильтры из формы
    getFilters(reportType) {
        const form = document.getElementById(`${reportType}ReportForm`);
        if (!form) return {};

        const filters = {};
        const formData = new FormData(form);

        for (const [key, value] of formData.entries()) {
            if (value) {
                filters[key] = value;
            }
        }

        return filters;
    },

    // Экспорт отчёта
    async exportReport(reportType) {
        try {
            window.app.showLoading();
            const filters = this.getFilters(reportType);
            const result = await Reports.exportToExcel(reportType, filters);

            if (result.success) {
                window.app.showNotification('Отчёт экспортирован', 'success');
            } else {
                window.app.showNotification(result.error || 'Ошибка экспорта', 'error');
            }
        } catch (error) {
            window.app.showNotification('Ошибка экспорта', 'error');
        } finally {
            window.app.hideLoading();
        }
    }
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Установка значений дат по умолчанию
    const startDateInputs = document.querySelectorAll('input[name="start_date"]');
    const endDateInputs = document.querySelectorAll('input[name="end_date"]');

    const today = new Date().toISOString().split('T')[0];
    const monthAgo = new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0];

    startDateInputs.forEach(input => {
        input.value = monthAgo;
    });

    endDateInputs.forEach(input => {
        input.value = today;
    });
});

// Экспорт для использования в других модулях
window.Reports = Reports;
window.ReportsUI = ReportsUI;