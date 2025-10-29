// Модуль админ панели

const Admin = {
    // Получить общую статистику
    async getStats() {
        try {
            const response = await API.get('/dashboard/stats');
            return response;
        } catch (error) {
            console.error('Error getting stats:', error);
            throw error;
        }
    },

    // Получить последние продажи
    async getRecentSales(limit = 10) {
        try {
            const response = await API.get(`/dashboard/recent-sales?limit=${limit}`);
            return response;
        } catch (error) {
            console.error('Error getting recent sales:', error);
            throw error;
        }
    },

    // Получить топ сотрудников
    async getTopEmployees(limit = 5) {
        try {
            const response = await API.get(`/dashboard/top-employees?limit=${limit}`);
            return response;
        } catch (error) {
            console.error('Error getting top employees:', error);
            throw error;
        }
    },

    // Получить все склады
    async getWarehouses() {
        try {
            const response = await API.get('/admin/warehouses');
            return response;
        } catch (error) {
            console.error('Error getting warehouses:', error);
            throw error;
        }
    },

    // Получить склад по ID
    async getWarehouse(id) {
        try {
            const response = await API.get(`/admin/warehouses/${id}`);
            return response;
        } catch (error) {
            console.error('Error getting warehouse:', error);
            throw error;
        }
    },

    // Создать склад
    async createWarehouse(data) {
        try {
            const response = await API.post('/admin/warehouses', data);
            return response;
        } catch (error) {
            console.error('Error creating warehouse:', error);
            throw error;
        }
    },

    // Обновить склад
    async updateWarehouse(id, data) {
        try {
            const response = await API.put(`/admin/warehouses/${id}`, data);
            return response;
        } catch (error) {
            console.error('Error updating warehouse:', error);
            throw error;
        }
    },

    // Удалить склад
    async deleteWarehouse(id) {
        try {
            const response = await API.delete(`/admin/warehouses/${id}`);
            return response;
        } catch (error) {
            console.error('Error deleting warehouse:', error);
            throw error;
        }
    }
};

// UI функции для админ панели
const AdminUI = {
    // Отобразить статистику на дашборде
    displayStats(stats) {
        const statsContainer = document.getElementById('statsContainer');
        if (!statsContainer) return;

        const statsHTML = `
            <div class="stat-card">
                <div class="stat-icon">🚗</div>
                <div class="stat-content">
                    <h3>Автомобили в наличии</h3>
                    <p class="stat-value">${stats.available_vehicles || 0}</p>
                </div>
            </div>

            <div class="stat-card">
                <div class="stat-icon">👥</div>
                <div class="stat-content">
                    <h3>Всего клиентов</h3>
                    <p class="stat-value">${stats.total_customers || 0}</p>
                </div>
            </div>

            <div class="stat-card">
                <div class="stat-icon">💰</div>
                <div class="stat-content">
                    <h3>Продаж за месяц</h3>
                    <p class="stat-value">${stats.sales_last_month || 0}</p>
                </div>
            </div>

            <div class="stat-card">
                <div class="stat-icon">📊</div>
                <div class="stat-content">
                    <h3>Общая выручка</h3>
                    <p class="stat-value">${window.app.formatCurrency(stats.total_revenue || 0)}</p>
                </div>
            </div>

            <div class="stat-card">
                <div class="stat-icon">👔</div>
                <div class="stat-content">
                    <h3>Сотрудников</h3>
                    <p class="stat-value">${stats.total_employees || 0}</p>
                </div>
            </div>

            <div class="stat-card">
                <div class="stat-icon">✅</div>
                <div class="stat-content">
                    <h3>Всего продаж</h3>
                    <p class="stat-value">${stats.total_sales || 0}</p>
                </div>
            </div>
        `;

        statsContainer.innerHTML = statsHTML;
    },

    // Отобразить последние продажи
    displayRecentSales(sales) {
        const container = document.getElementById('recentSalesContainer');
        if (!container) return;

        if (!sales || sales.length === 0) {
            container.innerHTML = '<div class="no-results">Нет последних продаж</div>';
            return;
        }

        const tableHTML = `
            <table class="table">
                <thead>
                    <tr>
                        <th>Дата</th>
                        <th>Автомобиль</th>
                        <th>Клиент</th>
                        <th>Сотрудник</th>
                        <th>Сумма</th>
                        <th>Статус</th>
                    </tr>
                </thead>
                <tbody>
                    ${sales.map(sale => `
                        <tr>
                            <td>${window.app.formatDate(sale.sale_date)}</td>
                            <td>${sale.model}</td>
                            <td>${sale.customer_name}</td>
                            <td>${sale.employee_name}</td>
                            <td>${window.app.formatCurrency(sale.final_price)}</td>
                            <td>
                                <span class="badge badge-${this.getStatusClass(sale.status)}">
                                    ${this.getStatusText(sale.status)}
                                </span>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Отобразить топ сотрудников
    displayTopEmployees(employees) {
        const container = document.getElementById('topEmployeesContainer');
        if (!container) return;

        if (!employees || employees.length === 0) {
            container.innerHTML = '<div class="no-results">Нет данных</div>';
            return;
        }

        const listHTML = `
            <div class="top-employees-list">
                ${employees.map((emp, index) => `
                    <div class="employee-item">
                        <div class="employee-rank">${index + 1}</div>
                        <div class="employee-info">
                            <h4>${emp.full_name}</h4>
                            <p>${emp.position}</p>
                        </div>
                        <div class="employee-stats">
                            <p class="sales-count">${emp.sales_count} продаж</p>
                            <p class="revenue">${window.app.formatCurrency(emp.total_revenue)}</p>
                        </div>
                    </div>
                `).join('')}
            </div>
        `;

        container.innerHTML = listHTML;
    },

    // Получить текст статуса
    getStatusText(status) {
        const statuses = {
            'pending': 'В ожидании',
            'completed': 'Завершена',
            'cancelled': 'Отменена'
        };
        return statuses[status] || status;
    },

    // Получить класс для статуса
    getStatusClass(status) {
        const classes = {
            'pending': 'warning',
            'completed': 'success',
            'cancelled': 'danger'
        };
        return classes[status] || 'primary';
    },

    // Инициализация дашборда
    async initDashboard() {
        try {
            window.app.showLoading();

            // Загружаем статистику
            const stats = await Admin.getStats();
            this.displayStats(stats);

            // Загружаем последние продажи
            const recentSales = await Admin.getRecentSales(10);
            this.displayRecentSales(recentSales);

            // Загружаем топ сотрудников
            const topEmployees = await Admin.getTopEmployees(5);
            this.displayTopEmployees(topEmployees);

            // Инициализируем графики
            if (typeof ChartsUI !== 'undefined') {
                ChartsUI.initSalesChart();
                ChartsUI.initVehicleStatusChart();
                ChartsUI.initEmployeeSalesChart();
            }

        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных дашборда', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Обновить дашборд
    async refreshDashboard() {
        await this.initDashboard();
        window.app.showNotification('Дашборд обновлён', 'success');
    }
};

// Sidebar navigation
const Sidebar = {
    init() {
        const sidebarLinks = document.querySelectorAll('.sidebar-menu a');
        const currentPath = window.location.pathname;

        sidebarLinks.forEach(link => {
            if (link.getAttribute('href') === currentPath) {
                link.classList.add('active');
            }

            link.addEventListener('click', function(e) {
                sidebarLinks.forEach(l => l.classList.remove('active'));
                this.classList.add('active');
            });
        });

        // Мобильное меню
        const menuToggle = document.getElementById('menuToggle');
        const sidebar = document.querySelector('.sidebar');

        if (menuToggle && sidebar) {
            menuToggle.addEventListener('click', () => {
                sidebar.classList.toggle('active');
            });

            // Закрытие при клике вне сайдбара
            document.addEventListener('click', (e) => {
                if (!sidebar.contains(e.target) && !menuToggle.contains(e.target)) {
                    sidebar.classList.remove('active');
                }
            });
        }
    }
};

// Инициализация при загрузке страницы админ панели
document.addEventListener('DOMContentLoaded', function() {
    // Проверка авторизации
    if (!Auth.isAuthenticated()) {
        window.location.href = '/login';
        return;
    }

    // Инициализация сайдбара
    Sidebar.init();

    // Инициализация дашборда если мы на странице дашборда
    if (window.location.pathname === '/admin/dashboard') {
        AdminUI.initDashboard();

        // Кнопка обновления
        const refreshBtn = document.getElementById('refreshDashboard');
        if (refreshBtn) {
            refreshBtn.addEventListener('click', () => {
                AdminUI.refreshDashboard();
            });
        }

        // Автообновление каждые 5 минут
        setInterval(() => {
            AdminUI.refreshDashboard();
        }, 5 * 60 * 1000);
    }

    // Инициализация страницы автомобилей
    if (window.location.pathname === '/admin/vehicles') {
        if (typeof VehiclesUI !== 'undefined') {
            VehiclesUI.loadVehiclesList();
        }

        const addVehicleBtn = document.getElementById('addVehicleBtn');
        if (addVehicleBtn) {
            addVehicleBtn.addEventListener('click', () => {
                VehiclesUI.openVehicleForm();
            });
        }
    }

    // Инициализация страницы продаж
    if (window.location.pathname === '/admin/sales') {
        if (typeof SalesUI !== 'undefined') {
            SalesUI.loadSalesList();
        }

        const addSaleBtn = document.getElementById('addSaleBtn');
        if (addSaleBtn) {
            addSaleBtn.addEventListener('click', () => {
                SalesUI.openSaleForm();
            });
        }
    }

    // Инициализация страницы клиентов
    if (window.location.pathname === '/admin/customers') {
        if (typeof CustomersUI !== 'undefined') {
            CustomersUI.loadCustomersList();
        }

        const addCustomerBtn = document.getElementById('addCustomerBtn');
        if (addCustomerBtn) {
            addCustomerBtn.addEventListener('click', () => {
                CustomersUI.openCustomerForm();
            });
        }
    }

    // Инициализация страницы отчётов
    if (window.location.pathname === '/admin/reports') {
        const reportType = document.getElementById('reportType');
        const generateReportBtn = document.getElementById('generateReportBtn');
        const exportReportBtn = document.getElementById('exportReportBtn');

        if (generateReportBtn) {
            generateReportBtn.addEventListener('click', () => {
                const type = reportType ? reportType.value : 'sales';
                ReportsUI.generateReport(type);
            });
        }

        if (exportReportBtn) {
            exportReportBtn.addEventListener('click', () => {
                const type = reportType ? reportType.value : 'sales';
                ReportsUI.exportReport(type);
            });
        }

        // Смена типа отчёта
        if (reportType) {
            reportType.addEventListener('change', (e) => {
                const selectedType = e.target.value;

                // Показываем/скрываем формы фильтров
                document.querySelectorAll('.report-filters').forEach(form => {
                    form.style.display = 'none';
                });

                const selectedForm = document.getElementById(`${selectedType}ReportForm`);
                if (selectedForm) {
                    selectedForm.style.display = 'block';
                }
            });
        }
    }

    // Отображение информации о пользователе
    const userInfoContainer = document.getElementById('userInfo');
    if (userInfoContainer) {
        Auth.getCurrentUser().then(user => {
            if (user) {
                userInfoContainer.innerHTML = `
                    <div class="user-info">
                        <span class="user-name">${user.full_name}</span>
                        <span class="user-role">${user.position}</span>
                    </div>
                `;
            }
        });
    }
});

// Функции для работы со складами
const WarehouseUI = {
    // Загрузить и отобразить склады
    async loadWarehouses() {
        console.log('Loading warehouses...');
        try {
            const response = await Admin.getWarehouses();
            console.log('Warehouses response:', response);
            if (response && response.data) {
                WarehouseUI.displayWarehouses(response.data);
            } else {
                console.error('No data in response:', response);
                document.getElementById('warehousesContainer').innerHTML = 
                    '<div class="error-state">Нет данных о складах</div>';
            }
        } catch (error) {
            console.error('Error loading warehouses:', error);
            document.getElementById('warehousesContainer').innerHTML = 
                '<div class="error-state">Ошибка загрузки складов: ' + error.message + '</div>';
        }
    },

    // Отобразить список складов
    displayWarehouses(warehouses) {
        console.log('Displaying warehouses:', warehouses);
        const container = document.getElementById('warehousesContainer');
        if (!container) {
            console.error('Warehouses container not found!');
            return;
        }

        if (!warehouses || warehouses.length === 0) {
            console.log('No warehouses found');
            container.innerHTML = '<div class="empty-state">Склады не найдены</div>';
            return;
        }

        const warehousesHTML = warehouses.map(warehouse => `
            <div class="warehouse-card">
                <div class="warehouse-header">
                    <h3>${warehouse.warehouse_name}</h3>
                    <div class="warehouse-actions">
                        <button class="btn btn-sm btn-primary" onclick="editWarehouse(${warehouse.warehouse_id})">
                            Редактировать
                        </button>
                        <button class="btn btn-sm btn-danger" onclick="deleteWarehouse(${warehouse.warehouse_id})">
                            Удалить
                        </button>
                    </div>
                </div>
                <div class="warehouse-info">
                    <div class="info-item">
                        <span class="label">Город:</span>
                        <span class="value">${warehouse.city}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">Вместимость:</span>
                        <span class="value">${warehouse.capacity} единиц</span>
                    </div>
                    <div class="info-item">
                        <span class="label">Техника в наличии:</span>
                        <span class="value">${warehouse.vehicles_in_stock || 0}</span>
                    </div>
                    <div class="info-item">
                        <span class="label">Статус:</span>
                        <span class="value ${warehouse.is_active ? 'active' : 'inactive'}">
                            ${warehouse.is_active ? 'Активен' : 'Неактивен'}
                        </span>
                    </div>
                </div>
            </div>
        `).join('');

        container.innerHTML = warehousesHTML;
    }
};

// Глобальные функции для работы со складами
window.editWarehouse = function(id) {
    console.log('Edit warehouse:', id);
    // TODO: Implement edit functionality
};

window.deleteWarehouse = function(id) {
    if (confirm('Вы уверены, что хотите удалить этот склад?')) {
        Admin.deleteWarehouse(id).then(() => {
            WarehouseUI.loadWarehouses();
        }).catch(error => {
            console.error('Error deleting warehouse:', error);
            alert('Ошибка при удалении склада');
        });
    }
};

// Инициализация страницы складов
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM loaded, checking for warehouses container...');
    // Проверяем, находимся ли мы на странице складов
    const container = document.getElementById('warehousesContainer');
    if (container) {
        console.log('Warehouses container found, loading warehouses...');
        WarehouseUI.loadWarehouses();
    } else {
        console.log('Warehouses container not found');
    }
});

// Экспорт для использования в других модулях
window.Admin = Admin;
window.AdminUI = AdminUI;
window.WarehouseUI = WarehouseUI;
window.Sidebar = Sidebar;