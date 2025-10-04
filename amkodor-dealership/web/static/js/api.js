// API Client для работы с backend
const API = {
    baseURL: '/api',

    // Получение токена из localStorage
    getToken() {
        return localStorage.getItem('token');
    },

    // Установка токена
    setToken(token) {
        localStorage.setItem('token', token);
    },

    // Удаление токена
    removeToken() {
        localStorage.removeItem('token');
    },

    // Базовый метод для запросов
    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const token = this.getToken();

        const config = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...options.headers,
            },
        };

        // Добавление токена для защищенных endpoints
        if (token && endpoint.includes('/admin')) {
            config.headers['Authorization'] = `Bearer ${token}`;
        }

        try {
            const response = await fetch(url, config);

            // Обработка ошибок авторизации
            if (response.status === 401) {
                this.removeToken();
                window.location.href = '/login';
                throw new Error('Unauthorized');
            }

            // Обработка других ошибок
            if (!response.ok) {
                const error = await response.json();
                throw new Error(error.error || 'Request failed');
            }

            const data = await response.json();
            return data.data || data;

        } catch (error) {
            console.error('API Request Error:', error);
            throw error;
        }
    },

    // GET запрос
    async get(endpoint, params = {}) {
        const queryString = new URLSearchParams(params).toString();
        const url = queryString ? `${endpoint}?${queryString}` : endpoint;
        return this.request(url, { method: 'GET' });
    },

    // POST запрос
    async post(endpoint, data) {
        return this.request(endpoint, {
            method: 'POST',
            body: JSON.stringify(data),
        });
    },

    // PUT запрос
    async put(endpoint, data) {
        return this.request(endpoint, {
            method: 'PUT',
            body: JSON.stringify(data),
        });
    },

    // DELETE запрос
    async delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    },

    // Auth endpoints
    auth: {
        async login(email, password) {
            const data = await API.post('/auth/login', { email, password });
            if (data.token) {
                API.setToken(data.token);
            }
            return data;
        },

        logout() {
            API.removeToken();
            window.location.href = '/login';
        },
    },

    // Vehicles endpoints
    vehicles: {
        async getAll(params = {}) {
            return API.get('/vehicles', params);
        },

        async getById(id) {
            return API.get(`/vehicles/${id}`);
        },

        async search(params) {
            return API.get('/vehicles/search', params);
        },

        async create(data) {
            return API.post('/admin/vehicles', data);
        },

        async update(id, data) {
            return API.put(`/admin/vehicles/${id}`, data);
        },

        async delete(id) {
            return API.delete(`/admin/vehicles/${id}`);
        },

        async getHistory(id) {
            return API.get(`/admin/vehicles/${id}/history`);
        },
    },

    // Customers endpoints
    customers: {
        async getAll(params = {}) {
            return API.get('/admin/customers', params);
        },

        async getById(id) {
            return API.get(`/admin/customers/${id}`);
        },

        async search(params) {
            return API.get('/admin/customers/search', params);
        },

        async create(data) {
            return API.post('/admin/customers', data);
        },

        async update(id, data) {
            return API.put(`/admin/customers/${id}`, data);
        },

        async delete(id) {
            return API.delete(`/admin/customers/${id}`);
        },
    },

    // Corporate Clients endpoints
    corporateClients: {
        async getAll(params = {}) {
            return API.get('/admin/corporate-clients', params);
        },

        async getById(id) {
            return API.get(`/admin/corporate-clients/${id}`);
        },

        async create(data) {
            return API.post('/admin/corporate-clients', data);
        },

        async update(id, data) {
            return API.put(`/admin/corporate-clients/${id}`, data);
        },

        async delete(id) {
            return API.delete(`/admin/corporate-clients/${id}`);
        },
    },

    // Sales endpoints
    sales: {
        async getAll(params = {}) {
            return API.get('/admin/sales', params);
        },

        async getById(id) {
            return API.get(`/admin/sales/${id}`);
        },

        async create(data) {
            return API.post('/admin/sales', data);
        },

        async update(id, data) {
            return API.put(`/admin/sales/${id}`, data);
        },

        async delete(id) {
            return API.delete(`/admin/sales/${id}`);
        },

        async getHistory(id) {
            return API.get(`/admin/sales/${id}/history`);
        },
    },

    // Employees endpoints
    employees: {
        async getAll(params = {}) {
            return API.get('/admin/employees', params);
        },

        async getById(id) {
            return API.get(`/admin/employees/${id}`);
        },

        async create(data) {
            return API.post('/admin/employees', data);
        },

        async update(id, data) {
            return API.put(`/admin/employees/${id}`, data);
        },

        async delete(id) {
            return API.delete(`/admin/employees/${id}`);
        },
    },

    // Warehouses endpoints
    warehouses: {
        async getAll(params = {}) {
            return API.get('/admin/warehouses', params);
        },

        async getById(id) {
            return API.get(`/admin/warehouses/${id}`);
        },

        async getStatistics(id) {
            return API.get(`/admin/warehouses/${id}/statistics`);
        },

        async create(data) {
            return API.post('/admin/warehouses', data);
        },

        async update(id, data) {
            return API.put(`/admin/warehouses/${id}`, data);
        },
    },

    // Service Orders endpoints
    service: {
        async getAllOrders(params = {}) {
            return API.get('/admin/service-orders', params);
        },

        async getOrderById(id) {
            return API.get(`/admin/service-orders/${id}`);
        },

        async createOrder(data) {
            return API.post('/admin/service-orders', data);
        },

        async updateOrder(id, data) {
            return API.put(`/admin/service-orders/${id}`, data);
        },

        async completeOrder(id) {
            return API.post(`/admin/service-orders/${id}/complete`);
        },

        async getAllParts(params = {}) {
            return API.get('/admin/spare-parts', params);
        },

        async getPartById(id) {
            return API.get(`/admin/spare-parts/${id}`);
        },

        async createPart(data) {
            return API.post('/admin/spare-parts', data);
        },

        async updatePart(id, data) {
            return API.put(`/admin/spare-parts/${id}`, data);
        },
    },

    // Test Drives endpoints
    testDrives: {
        async getAll(params = {}) {
            return API.get('/admin/test-drives', params);
        },

        async create(data) {
            return API.post('/test-drives', data);
        },

        async update(id, data) {
            return API.put(`/admin/test-drives/${id}`, data);
        },
    },

    // Reports endpoints
    reports: {
        async sales(params) {
            return API.get('/admin/reports/sales', params);
        },

        async inventory(params) {
            return API.get('/admin/reports/inventory', params);
        },

        async exportSales(params) {
            const queryString = new URLSearchParams(params).toString();
            const url = `/admin/reports/export/sales?${queryString}`;
            window.open(url, '_blank');
        },

        async exportInventory(params) {
            const queryString = new URLSearchParams(params).toString();
            const url = `/admin/reports/export/inventory?${queryString}`;
            window.open(url, '_blank');
        },
    },

    // Dashboard endpoints
    dashboard: {
        async getStatistics() {
            return API.get('/admin/dashboard');
        },

        async getChartData() {
            return API.get('/admin/dashboard/charts');
        },
    },
};

// Утилиты для работы с датами
const DateUtils = {
    format(date, format = 'YYYY-MM-DD') {
        const d = new Date(date);
        const year = d.getFullYear();
        const month = String(d.getMonth() + 1).padStart(2, '0');
        const day = String(d.getDate()).padStart(2, '0');

        if (format === 'YYYY-MM-DD') {
            return `${year}-${month}-${day}`;
        }
        if (format === 'DD.MM.YYYY') {
            return `${day}.${month}.${year}`;
        }
        return date;
    },

    formatDateTime(date) {
        const d = new Date(date);
        return d.toLocaleString('ru-BY');
    },
};

// Утилиты для работы с числами
const NumberUtils = {
    formatCurrency(value) {
        return new Intl.NumberFormat('ru-BY', {
            style: 'decimal',
            minimumFractionDigits: 2,
            maximumFractionDigits: 2
        }).format(value);
    },

    formatNumber(value) {
        return new Intl.NumberFormat('ru-BY').format(value);
    },
};

// Проверка авторизации при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    const token = API.getToken();
    const isAdminPage = window.location.pathname.startsWith('/admin');
    const isLoginPage = window.location.pathname === '/login';

    if (isAdminPage && !token) {
        window.location.href = '/login';
    }

    if (isLoginPage && token) {
        window.location.href = '/admin/dashboard';
    }
});