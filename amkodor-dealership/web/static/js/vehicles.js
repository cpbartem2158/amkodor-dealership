// Модуль работы с автомобилями

const Vehicles = {
    // Получить все автомобили
    async getAll() {
        try {
            const response = await API.get('/vehicles');
            return response;
        } catch (error) {
            console.error('Error getting vehicles:', error);
            throw error;
        }
    },

    // Получить автомобиль по ID
    async getById(id) {
        try {
            const response = await API.get(`/vehicles/${id}`);
            return response;
        } catch (error) {
            console.error('Error getting vehicle:', error);
            throw error;
        }
    },

    // Создать автомобиль
    async create(vehicleData) {
        try {
            const response = await API.post('/vehicles', vehicleData);
            return response;
        } catch (error) {
            console.error('Error creating vehicle:', error);
            throw error;
        }
    },

    // Обновить автомобиль
    async update(id, vehicleData) {
        try {
            const response = await API.put(`/vehicles/${id}`, vehicleData);
            return response;
        } catch (error) {
            console.error('Error updating vehicle:', error);
            throw error;
        }
    },

    // Удалить автомобиль
    async delete(id) {
        try {
            const response = await API.delete(`/vehicles/${id}`);
            return response;
        } catch (error) {
            console.error('Error deleting vehicle:', error);
            throw error;
        }
    },

    // Поиск автомобилей
    async search(filters) {
        try {
            const queryString = new URLSearchParams(filters).toString();
            const response = await API.get(`/vehicles/search?${queryString}`);
            return response;
        } catch (error) {
            console.error('Error searching vehicles:', error);
            throw error;
        }
    },

    // Получить категории
    async getCategories() {
        try {
            const response = await API.get('/categories');
            return response;
        } catch (error) {
            console.error('Error getting categories:', error);
            throw error;
        }
    }
};

// UI функции для работы с автомобилями
const VehiclesUI = {
    // Отобразить список автомобилей
    displayVehicles(vehicles, containerId = 'vehiclesContainer') {
        const container = document.getElementById(containerId);
        if (!container) return;

        if (!vehicles || vehicles.length === 0) {
            container.innerHTML = '<div class="no-results">Автомобили не найдены</div>';
            return;
        }

        container.innerHTML = vehicles.map(vehicle => `
            <div class="vehicle-card" data-id="${vehicle.id}">
                <div class="vehicle-image"></div>
                <div class="vehicle-info">
                    <h3 class="vehicle-title">${vehicle.model}</h3>
                    <span class="badge badge-${this.getStatusClass(vehicle.status)}">
                        ${this.getStatusText(vehicle.status)}
                    </span>
                    <div class="vehicle-details">
                        <div class="detail-row">
                            <span class="detail-label">Год:</span>
                            <span>${vehicle.year}</span>
                        </div>
                        <div class="detail-row">
                            <span class="detail-label">VIN:</span>
                            <span>${vehicle.vin}</span>
                        </div>
                        <div class="detail-row">
                            <span class="detail-label">Цвет:</span>
                            <span>${vehicle.color}</span>
                        </div>
                        <div class="detail-row">
                            <span class="detail-label">Пробег:</span>
                            <span>${window.app.formatNumber(vehicle.mileage)} км</span>
                        </div>
                    </div>
                    <div class="vehicle-price">${window.app.formatCurrency(vehicle.price)}</div>
                    <button class="btn btn-primary" onclick="VehiclesUI.viewDetails(${vehicle.id})">
                        Подробнее
                    </button>
                </div>
            </div>
        `).join('');
    },

    // Просмотр деталей автомобиля
    viewDetails(id) {
        window.location.href = `/vehicle/${id}`;
    },

    // Получить текст статуса
    getStatusText(status) {
        const statuses = {
            'available': 'В наличии',
            'reserved': 'Зарезервирован',
            'sold': 'Продан',
            'in_service': 'На обслуживании'
        };
        return statuses[status] || status;
    },

    // Получить класс для статуса
    getStatusClass(status) {
        const classes = {
            'available': 'success',
            'reserved': 'warning',
            'sold': 'danger',
            'in_service': 'info'
        };
        return classes[status] || 'primary';
    },

    // Форма создания/редактирования автомобиля
    openVehicleForm(vehicleId = null) {
        const modal = document.getElementById('vehicleModal');
        const form = document.getElementById('vehicleForm');
        const title = document.getElementById('modalTitle');

        if (vehicleId) {
            title.textContent = 'Редактировать автомобиль';
            this.loadVehicleData(vehicleId, form);
        } else {
            title.textContent = 'Добавить автомобиль';
            form.reset();
        }

        form.dataset.vehicleId = vehicleId || '';
        window.app.openModal('vehicleModal');
    },

    // Загрузка данных автомобиля для редактирования
    async loadVehicleData(vehicleId, form) {
        try {
            window.app.showLoading();
            const vehicle = await Vehicles.getById(vehicleId);

            if (vehicle) {
                form.querySelector('[name="vin"]').value = vehicle.vin;
                form.querySelector('[name="model"]').value = vehicle.model;
                form.querySelector('[name="year"]').value = vehicle.year;
                form.querySelector('[name="color"]').value = vehicle.color;
                form.querySelector('[name="mileage"]').value = vehicle.mileage;
                form.querySelector('[name="price"]').value = vehicle.price;
                form.querySelector('[name="status"]').value = vehicle.status;
                form.querySelector('[name="category_id"]').value = vehicle.category_id;
            }
        } catch (error) {
            window.app.showNotification('Ошибка загрузки данных', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Сохранить автомобиль
    async saveVehicle(formData, vehicleId) {
        try {
            window.app.showLoading();

            let result;
            if (vehicleId) {
                result = await Vehicles.update(vehicleId, formData);
            } else {
                result = await Vehicles.create(formData);
            }

            if (result.success) {
                window.app.showNotification(
                    vehicleId ? 'Автомобиль обновлён' : 'Автомобиль добавлен',
                    'success'
                );
                window.app.closeModal('vehicleModal');
                this.loadVehiclesList();
            } else {
                window.app.showNotification(result.error || 'Ошибка сохранения', 'error');
            }
        } catch (error) {
            window.app.showNotification('Ошибка сохранения', 'error');
        } finally {
            window.app.hideLoading();
        }
    },

    // Удалить автомобиль
    async deleteVehicle(vehicleId) {
        window.app.confirmAction('Вы уверены, что хотите удалить этот автомобиль?', async () => {
            try {
                window.app.showLoading();
                const result = await Vehicles.delete(vehicleId);

                if (result.success) {
                    window.app.showNotification('Автомобиль удалён', 'success');
                    this.loadVehiclesList();
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

    // Загрузить список автомобилей
    async loadVehiclesList() {
        try {
            window.app.showLoading();
            const vehicles = await Vehicles.getAll();
            this.displayVehicles(vehicles);
        } catch (error) {
            window.app.showNotification('Ошибка загрузки автомобилей', 'error');
        } finally {
            window.app.hideLoading();
        }
    }
};

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    // Обработчик формы автомобиля
    const vehicleForm = document.getElementById('vehicleForm');
    if (vehicleForm) {
        vehicleForm.addEventListener('submit', async function(e) {
            e.preventDefault();

            const formData = {
                vin: this.querySelector('[name="vin"]').value,
                model: this.querySelector('[name="model"]').value,
                year: parseInt(this.querySelector('[name="year"]').value),
                color: this.querySelector('[name="color"]').value,
                mileage: parseInt(this.querySelector('[name="mileage"]').value),
                price: parseFloat(this.querySelector('[name="price"]').value),
                status: this.querySelector('[name="status"]').value,
                category_id: parseInt(this.querySelector('[name="category_id"]').value)
            };

            const vehicleId = this.dataset.vehicleId;
            await VehiclesUI.saveVehicle(formData, vehicleId);
        });
    }

    // Загрузка категорий для селекта
    const categorySelect = document.querySelector('[name="category_id"]');
    if (categorySelect) {
        Vehicles.getCategories().then(categories => {
            categorySelect.innerHTML = categories.map(cat =>
                `<option value="${cat.id}">${cat.name}</option>`
            ).join('');
        });
    }
});

// Экспорт для использования в других модулях
window.Vehicles = Vehicles;
window.VehiclesUI = VehiclesUI;