// Модуль работы с продажами

const Sales = {
    // Получить все продажи
    async getAll() {
        try {
            const response = await API.get('/sales');
            return response;
        } catch (error) {
            console.error('Error getting sales:', error);
            throw error;
        }
    },

    // Получить продажу по ID
    async getById(id) {
        try {
            const response = await API.get(`/sales/${id}`);
            return response;
        } catch (error) {
            console.error('Error getting sale:', error);
            throw error;
        }
    },

    // Создать продажу
    async create(saleData) {
        try {
            const response = await API.post('/sales', saleData);
            return response;
        } catch (error) {
            console.error('Error creating sale:', error);
            throw error;
        }
    },

    // Обновить продажу
    async update(id, saleData) {
        try {
            const response = await API.put(`/sales/${id}`, saleData);
            return response;
        } catch (error) {
            console.error('Error updating sale:', error);
            throw error;
        }
    },

    // Отменить продажу
    async cancel(id) {
        try {
            const response = await API.post(`/sales/${id}/cancel`);
            return response;
        } catch (error) {
            console.error('Error cancelling sale:', error);
            throw error;
        }
    }
};

// UI функции для работы с продажами
const SalesUI = {
    // Отобразить список продаж
    displaySales(sales, containerId = 'salesContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!sales || sales.length === 0) {
            container.innerHTML = '<div class="no-results">Продажи не найдены</div>';
            return;
        }

        const tableHTML = `
            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>Дата</th>
                        <th>Автомобиль</th>
                        <th>Клиент</th>
                        <th>Сотрудник</th>
                        <th>Сумма</th>
                        <th>Скидка</th>
                        <th>Статус</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${sales.map(sale => `
                        <tr>
                            <td>${sale.id}</td>
                            <td>${window.app.formatDate(sale.sale_date)}</td>
                            <td>${sale.vehicle_model} (${sale.vehicle_vin})</td>
                            <td>${sale.customer_name}</td>
                            <td>${sale.employee_name}</td>
                            <td>${window.app.formatCurrency(sale.final_price)}</td>
                            <td>${window.app.formatCurrency(sale.discount)}</td>
                            <td>
                                <span class="badge badge-${this.getStatusClass(sale.status)}">
                                    ${this.getStatusText(sale.status)}
                                </span>
                            </td>
                            <td>
                                <button class="btn btn-sm" onclick="SalesUI.viewDetails(${sale.id})">
                                    Просмотр
                                </button>
                                ${sale.status === 'completed' ? `
                                    <button class="btn btn-sm btn-danger" onclick="SalesUI.cancelSale(${sale.id})">
                                        Отменить
                                    </button>
                                ` : ''}
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Просмотр деталей продажи
    async viewDetails(id) {
        try {
            window.app.showLoading();
            const sale = await Sales.getById(id);

            if (sale) {
                this.showSaleDetails(sale);
            }
        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Показать детали продажи в модальном окне
    showSaleDetails(sale) {
        const modalContent = `
            <h2>Продажа #${sale.id}</h2>
            <div class="sale-details">
                <div class="detail-section">
                    <h3>Информация о продаже</h3>
                    <p><strong>Дата:</strong> ${window.app.formatDateTime(sale.sale_date)}</p>
                    <p><strong>Статус:</strong> 
                        <span class="badge badge-${this.getStatusClass(sale.status)}">
                            ${this.getStatusText(sale.status)}
                        </span>
                    </p>
                </div>
                <div class="detail-section">
                    <h3>Автомобиль</h3>
                    <p><strong>Модель:</strong> ${sale.vehicle_model}</p>
                    <p><strong>VIN:</strong> ${sale.vehicle_vin}</p>
                    <p><strong>Год:</strong> ${sale.vehicle_year}</p>
                </div>
                <div class="detail-section">
                    <h3>Клиент</h3>
                    <p><strong>Имя:</strong> ${sale.customer_name}</p>
                    <p><strong>Телефон:</strong> ${sale.customer_phone || 'Не указан'}</p>
                    <p><strong>Email:</strong> ${sale.customer_email || 'Не указан'}</p>
                </div>
                <div class="detail-section">
                    <h3>Сотрудник</h3>
                    <p><strong>Имя:</strong> ${sale.employee_name}</p>
                </div>
                <div class="detail-section">
                    <h3>Финансы</h3>
                    <p><strong>Скидка:</strong> ${window.app.formatCurrency(sale.discount)}</p>
                    <p><strong>Итоговая сумма:</strong> ${window.app.formatCurrency(sale.final_price)}</p>
                </div>
            </div>
        `;

        const modal = document.getElementById('saleDetailsModal');
        if (modal) {
            modal.querySelector('.modal-content').innerHTML = modalContent;
            window.app.openModal('saleDetailsModal');
        }
    },

    // Форма создания продажи
    openSaleForm() {
        window.app.openModal('saleModal');
        this.loadFormData();
    },

    // Загрузка данных для формы
    async loadFormData() {
        try {
            // Загружаем доступные автомобили
            const vehicles = await Vehicles.getAll();
            const availableVehicles = vehicles.filter(v => v.status === 'available');

            const vehicleSelect = document.querySelector('[name="vehicle_id"]');
            if (vehicleSelect) {
                vehicleSelect.innerHTML = '<option value="">Выберите автомобиль</option>' +
                    availableVehicles.map(v =>
                        `<option value="${v.id}">${v.model} - ${v.vin}</option>`
                    ).join('');
            }

            // Загружаем клиентов
            const customers = await API.get('/customers');
            const customerSelect = document.querySelector('[name="customer_id"]');
            if (customerSelect) {
                customerSelect.innerHTML = '<option value="">Выберите клиента</option>' +
                    customers.map(c =>
                        `<option value="${c.id}">${c.full_name}</option>`
                    ).join('');
            }

            // Загружаем сотрудников
            const employees = await API.get('/employees');
            const employeeSelect = document.querySelector('[name="employee_id"]');
            if (employeeSelect) {
                employeeSelect.innerHTML = '<option value="">Выберите сотрудника</option>' +
                    employees.map(e =>
                        `<option value="${e.id}">${e.full_name}</option>`
                    ).join('');
            }
        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных', 'error');
        }
    },

    // Сохранить продажу
    async saveSale(formData) {
        try {
            window.app.showLoading();
            const result = await Sales.create(formData);

            if (result.success) {
                window.app.showNotification('Продажа создана', 'success');
                window.app.closeModal('saleModal');
                this.loadSalesList();
            } else {
                window.app.showNotification(result.error || 'Ошибка создания', 'error');
            }
        } catch (error) {
            window.app.showNotification('Ошибка создания', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Отменить продажу
    async cancelSale(saleId) {
        window.app.confirmAction('Вы уверены, что хотите отменить эту продажу?', async () => {
            try {
                window.app.showLoading();
                const result = await Sales.cancel(saleId);

                if (result.success) {
                    window.app.showNotification('Продажа отменена', 'success');
                    this.loadSalesList();
                } else {
                    window.app.showNotification(result.error || 'Ошибка отмены', 'error');
                }
            } catch (error) {
                window.app.showNotification('Ошибка отмены', 'error');
            } finally {
                window.app.hideLoading();
            }
        });
    },

    // Загрузить список продаж
    async loadSalesList() {
        try {
            window.app.showLoading();
            const sales = await Sales.getAll();
            this.displaySales(sales);
        } catch (error) {
            window.app.showNotification('Ошибка загрузки продаж', 'error');
        } finally {
            window.app.hideLoading();
        }
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
    }
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Обработчик формы продажи
    const saleForm = document.getElementById('saleForm');
    if (saleForm) {
        saleForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            const formData = {
                vehicle_id: parseInt(this.querySelector('[name="vehicle_id"]').value),
                customer_id: parseInt(this.querySelector('[name="customer_id"]').value),
                employee_id: parseInt(this.querySelector('[name="employee_id"]').value),
                sale_date: this.querySelector('[name="sale_date"]').value,
                discount: parseFloat(this.querySelector('[name="discount"]').value) || 0,
                final_price: parseFloat(this.querySelector('[name="final_price"]').value),
                status: 'completed'
            };

            await SalesUI.saveSale(formData);
        });

        // Автоматический расчёт финальной цены
        const vehicleSelect = saleForm.querySelector('[name="vehicle_id"]');
        const discountInput = saleForm.querySelector('[name="discount"]');
        const finalPriceInput = saleForm.querySelector('[name="final_price"]');

        if (vehicleSelect && discountInput && finalPriceInput) {
            const calculateFinalPrice = async () => {
                const vehicleId = vehicleSelect.value;
                const discount = parseFloat(discountInput.value) || 0;

                if (vehicleId) {
                    const vehicle = await Vehicles.getById(vehicleId);
                    if (vehicle) {
                        const finalPrice = vehicle.price - discount;
                        finalPriceInput.value = finalPrice.toFixed(2);
                    }
                }
            };

            vehicleSelect.addEventListener('change', calculateFinalPrice);
            discountInput.addEventListener('input', calculateFinalPrice);
        }
    }
});

// Экспорт для использования в других модулях
window.Sales = Sales;
window.SalesUI = SalesUI;