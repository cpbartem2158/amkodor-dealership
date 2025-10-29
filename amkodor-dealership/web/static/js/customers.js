// Модуль работы с клиентами

const Customers = {
    // Получить всех клиентов
    async getAll() {
        try {
            const response = await API.get('/customers');
            return response;
        } catch (error) {
            console.error('Error getting customers:', error);
            throw error;
        }
    },

    // Получить клиента по ID
    async getById(id) {
        try {
            const response = await API.get(`/customers/${id}`);
            return response;
        } catch (error) {
            console.error('Error getting customer:', error);
            throw error;
        }
    },

    // Создать клиента
    async create(customerData) {
        try {
            const response = await API.post('/customers', customerData);
            return response;
        } catch (error) {
            console.error('Error creating customer:', error);
            throw error;
        }
    },

    // Обновить клиента
    async update(id, customerData) {
        try {
            const response = await API.put(`/customers/${id}`, customerData);
            return response;
        } catch (error) {
            console.error('Error updating customer:', error);
            throw error;
        }
    },

    // Удалить клиента
    async delete(id) {
        try {
            const response = await API.delete(`/customers/${id}`);
            return response;
        } catch (error) {
            console.error('Error deleting customer:', error);
            throw error;
        }
    },

    // Поиск клиентов
    async search(filters) {
        try {
            const queryString = new URLSearchParams(filters).toString();
            const response = await API.get(`/customers/search?${queryString}`);
            return response;
        } catch (error) {
            console.error('Error searching customers:', error);
            throw error;
        }
    }
};

// UI функции для работы с клиентами
const CustomersUI = {
    // Отобразить список клиентов
    displayCustomers(customers, containerId = 'customersContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!customers || customers.length === 0) {
            container.innerHTML = '<div class="no-results">Клиенты не найдены</div>';
            return;
        }

        const tableHTML = `
            <table class="table">
                <thead>
                    <tr>
                        <th>ID</th>
                        <th>ФИО</th>
                        <th>Телефон</th>
                        <th>Email</th>
                        <th>Адрес</th>
                        <th>Покупок</th>
                        <th>Потрачено</th>
                        <th>Действия</th>
                    </tr>
                </thead>
                <tbody>
                    ${customers.map(customer => `
                        <tr>
                            <td>${customer.id}</td>
                            <td>${customer.full_name}</td>
                            <td>${customer.phone}</td>
                            <td>${customer.email || '-'}</td>
                            <td>${customer.address || '-'}</td>
                            <td>${customer.purchases_count || 0}</td>
                            <td>${window.app.formatCurrency(customer.total_spent || 0)}</td>
                            <td>
                                <button class="btn btn-sm" onclick="CustomersUI.viewDetails(${customer.id})">
                                    Просмотр
                                </button>
                                <button class="btn btn-sm" onclick="CustomersUI.openCustomerForm(${customer.id})">
                                    Изменить
                                </button>
                                <button class="btn btn-sm btn-danger" onclick="CustomersUI.deleteCustomer(${customer.id})">
                                    Удалить
                                </button>
                            </td>
                        </tr>
                    `).join('')}
                </tbody>
            </table>
        `;

        container.innerHTML = tableHTML;
    },

    // Просмотр деталей клиента
    async viewDetails(id) {
        try {
            window.app.showLoading();
            const customer = await Customers.getById(id);

            if (customer) {
                this.showCustomerDetails(customer);
            }
        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Показать детали клиента
    showCustomerDetails(customer) {
        const modalContent = `
            <h2>${customer.full_name}</h2>
            <div class="customer-details">
                <div class="detail-section">
                    <h3>Контактная информация</h3>
                    <p><strong>Телефон:</strong> ${customer.phone}</p>
                    <p><strong>Email:</strong> ${customer.email || 'Не указан'}</p>
                    <p><strong>Адрес:</strong> ${customer.address || 'Не указан'}</p>
                    <p><strong>Паспорт:</strong> ${customer.passport_number || 'Не указан'}</p>
                </div>
                <div class="detail-section">
                    <h3>Статистика</h3>
                    <p><strong>Количество покупок:</strong> ${customer.purchases_count || 0}</p>
                    <p><strong>Общая сумма покупок:</strong> ${window.app.formatCurrency(customer.total_spent || 0)}</p>
                    <p><strong>Дата регистрации:</strong> ${window.app.formatDate(customer.created_at)}</p>
                </div>
            </div>
        `;

        const modal = document.getElementById('customerDetailsModal');
        if (modal) {
            modal.querySelector('.modal-content').innerHTML = modalContent;
            window.app.openModal('customerDetailsModal');
        }
    },

    // Форма создания/редактирования клиента
    openCustomerForm(customerId = null) {
        const modal = document.getElementById('customerModal');
        const form = document.getElementById('customerForm');
        const title = document.getElementById('modalTitle');

        if (customerId) {
            title.textContent = 'Редактировать клиента';
            this.loadCustomerData(customerId, form);
        } else {
            title.textContent = 'Добавить клиента';
            form.reset();
        }

        form.dataset.customerId = customerId || '';
        window.app.openModal('customerModal');
    },

    // Загрузка данных клиента
    async loadCustomerData(customerId, form) {
        try {
            window.app.showLoading();
            const customer = await Customers.getById(customerId);

            if (customer) {
                form.querySelector('[name="full_name"]').value = customer.full_name;
                form.querySelector('[name="phone"]').value = customer.phone;
                form.querySelector('[name="email"]').value = customer.email || '';
                form.querySelector('[name="address"]').value = customer.address || '';
                form.querySelector('[name="passport_number"]').value = customer.passport_number || '';
            }
        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Сохранить клиента
    async saveCustomer(formData, customerId) {
        try {
            window.app.showLoading();

            let result;
            if (customerId) {
                result = await Customers.update(customerId, formData);
            } else {
                result = await Customers.create(formData);
            }

            if (result.success) {
                window.app.showNotification(
                    customerId ? 'Клиент обновлён' : 'Клиент добавлен',
                    'success'
                );
                window.app.closeModal('customerModal');
                this.loadCustomersList();
            } else {
                window.app.showNotification(result.error || 'Ошибка сохранения', 'error');
            }
        } catch (error) {
            window.app.showNotification('Ошибка сохранения', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Удалить клиента
    async deleteCustomer(customerId) {
        window.app.confirmAction('Вы уверены, что хотите удалить этого клиента?', async () => {
            try {
                window.app.showLoading();
                const result = await Customers.delete(customerId);

                if (result.success) {
                    window.app.showNotification('Клиент удалён', 'success');
                    this.loadCustomersList();
                } else {
                    window.app.showNotification(result.error || 'Ошибка удаления', 'error');
                }
            } catch (error) {
                window.app.showNotification('Ошибка удаления', 'error');
            } finally {
                window.app.hideLoading();
            }
        });
    },

    // Загрузить список клиентов
    async loadCustomersList() {
        try {
            window.app.showLoading();
            const customers = await Customers.getAll();
            this.displayCustomers(customers);
        } catch (error) {
            window.app.showNotification('Ошибка загрузки клиентов', 'error');
        } finally {
            window.app.hideLoading();
        }
    }
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Обработчик формы клиента
    const customerForm = document.getElementById('customerForm');
    if (customerForm) {
        customerForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            const formData = {
                full_name: this.querySelector('[name="full_name"]').value,
                phone: this.querySelector('[name="phone"]').value,
                email: this.querySelector('[name="email"]').value,
                address: this.querySelector('[name="address"]').value,
                passport_number: this.querySelector('[name="passport_number"]').value
            };

            const customerId = this.dataset.customerId;
            await CustomersUI.saveCustomer(formData, customerId);
        });
    }

    // Обработчик поиска
    const searchForm = document.getElementById('customerSearchForm');
    if (searchForm) {
        const searchInput = searchForm.querySelector('input[name="search"]');

        if (searchInput) {
            searchInput.addEventListener('input', window.app.debounce(async function() {
                const searchTerm = this.value;
                if (searchTerm.length >= 2) {
                    const results = await Customers.search({ name: searchTerm });
                    CustomersUI.displayCustomers(results);
                } else if (searchTerm.length === 0) {
                    CustomersUI.loadCustomersList();
                }
            }, 500));
        }
    }
});

// Экспорт для использования в других модулях
window.Customers = Customers;
window.CustomersUI = CustomersUI;