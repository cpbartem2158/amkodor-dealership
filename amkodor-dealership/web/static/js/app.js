// Главный JavaScript файл приложения

// Инициализация при загрузке страницы
document.addEventListener('DOMContentLoaded', function() {
    console.log('App initialized');
    initNavigation();
    initModals();
    initForms();
    checkAuth();
    testAPI();
});

// Навигация
function initNavigation() {
    const currentPath = window.location.pathname;
    const navLinks = document.querySelectorAll('.nav-links a');

    navLinks.forEach(link => {
        if (link.getAttribute('href') === currentPath) {
            link.classList.add('active');
        }
    });

    // Мобильное меню
    const menuToggle = document.querySelector('.menu-toggle');
    const navLinksContainer = document.querySelector('.nav-links');

    if (menuToggle) {
        menuToggle.addEventListener('click', () => {
            navLinksContainer.classList.toggle('active');
        });
    }
}

// Модальные окна
function initModals() {
    const modals = document.querySelectorAll('.modal');

    modals.forEach(modal => {
        const closeBtn = modal.querySelector('.modal-close');

        if (closeBtn) {
            closeBtn.addEventListener('click', () => {
                closeModal(modal);
            });
        }

        // Закрытие при клике вне модального окна
        modal.addEventListener('click', (e) => {
            if (e.target === modal) {
                closeModal(modal);
            }
        });
    });
}

function openModal(modalId) {
    const modal = document.getElementById(modalId);
    if (modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

function closeModal(modal) {
    if (typeof modal === 'string') {
        modal = document.getElementById(modal);
    }
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = '';
    }
}

// Формы
function initForms() {
    const forms = document.querySelectorAll('form[data-validate]');

    forms.forEach(form => {
        form.addEventListener('submit', (e) => {
            if (!validateForm(form)) {
                e.preventDefault();
            }
        });
    });
}

function validateForm(form) {
    let isValid = true;
    const inputs = form.querySelectorAll('[required]');

    inputs.forEach(input => {
        if (!input.value.trim()) {
            showError(input, 'Это поле обязательно');
            isValid = false;
        } else {
            clearError(input);
        }

        // Валидация email
        if (input.type === 'email' && input.value) {
            const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
            if (!emailRegex.test(input.value)) {
                showError(input, 'Некорректный email');
                isValid = false;
            }
        }

        // Валидация телефона
        if (input.type === 'tel' && input.value) {
            const phoneRegex = /^\+?[0-9]{10,15}$/;
            if (!phoneRegex.test(input.value.replace(/\s/g, ''))) {
                showError(input, 'Некорректный номер телефона');
                isValid = false;
            }
        }
    });

    return isValid;
}

function showError(input, message) {
    input.classList.add('error');

    let errorDiv = input.parentElement.querySelector('.form-error');
    if (!errorDiv) {
        errorDiv = document.createElement('div');
        errorDiv.className = 'form-error';
        input.parentElement.appendChild(errorDiv);
    }
    errorDiv.textContent = message;
}

function clearError(input) {
    input.classList.remove('error');
    const errorDiv = input.parentElement.querySelector('.form-error');
    if (errorDiv) {
        errorDiv.remove();
    }
}

// Уведомления
function showNotification(message, type = 'info') {
    const notification = document.createElement('div');
    notification.className = `alert alert-${type}`;
    notification.textContent = message;

    const container = document.querySelector('.notifications') || document.body;
    container.appendChild(notification);

    setTimeout(() => {
        notification.style.opacity = '0';
        setTimeout(() => notification.remove(), 300);
    }, 3000);
}

// Загрузчик
function showLoading() {
    const overlay = document.createElement('div');
    overlay.className = 'loading-overlay';
    overlay.innerHTML = '<div class="spinner"></div>';
    overlay.id = 'loadingOverlay';
    document.body.appendChild(overlay);
}

function hideLoading() {
    const overlay = document.getElementById('loadingOverlay');
    if (overlay) {
        overlay.remove();
    }
}

// Проверка авторизации
function checkAuth() {
    const token = localStorage.getItem('token');
    const userRole = localStorage.getItem('userRole');
    const protectedPages = ['/admin'];
    const currentPath = window.location.pathname;

    // Если пользователь на защищенной странице без токена - перенаправляем на логин
    if (protectedPages.some(page => currentPath.startsWith(page)) && !token) {
        window.location.href = '/login';
        return;
    }

    // Если пользователь на главной странице с токеном - проверяем роль
    if (currentPath === '/' && token && userRole === 'admin') {
        // Админы автоматически перенаправляются в админ панель
        window.location.href = '/admin/dashboard';
        return;
    }

    if (currentPath === '/' && token && userRole === 'user') {
        // Обычные пользователи перенаправляются в дашборд
        window.location.href = '/dashboard';
        return;
    }

    // Показываем кнопку "Выйти" для авторизованных пользователей на главной странице
    if (currentPath === '/' && token) {
        const loginBtn = document.getElementById('login-btn');
        const logoutBtn = document.getElementById('logout-btn');
        if (loginBtn) loginBtn.style.display = 'none';
        if (logoutBtn) logoutBtn.style.display = 'block';
    }
}

// Функция выхода из системы
function logout() {
    localStorage.removeItem('token');
    localStorage.removeItem('userRole');
    localStorage.removeItem('userName');
    window.location.href = '/';
}

// Форматирование чисел
function formatNumber(num) {
    return new Intl.NumberFormat('ru-RU').format(num);
}

function formatCurrency(amount) {
    return new Intl.NumberFormat('ru-RU', {
        style: 'currency',
        currency: 'BYN'
    }).format(amount);
}

// Форматирование дат
function formatDate(date) {
    return new Intl.DateTimeFormat('ru-RU').format(new Date(date));
}

// Тестирование API
async function testAPI() {
    try {
        console.log('Testing API connection...');
        const response = await fetch('/api/health');
        const data = await response.json();
        console.log('API Health:', data);
        
        // Тестируем статистику
        const statsResponse = await fetch('/api/stats');
        const stats = await statsResponse.json();
        console.log('API Stats:', stats);
        
        // Обновляем статистику на странице, если есть элементы
        updateStatsOnPage(stats.data);
        
    } catch (error) {
        console.error('API Test Error:', error);
    }
}

// Обновление статистики на странице
function updateStatsOnPage(stats) {
    // Ищем элементы статистики и обновляем их
    const totalVehicles = document.querySelector('[data-stat="total-vehicles"]');
    const availableVehicles = document.querySelector('[data-stat="available-vehicles"]');
    const totalSales = document.querySelector('[data-stat="total-sales"]');
    const totalRevenue = document.querySelector('[data-stat="total-revenue"]');
    const totalCustomers = document.querySelector('[data-stat="total-customers"]');
    
    if (totalVehicles) totalVehicles.textContent = stats.total_vehicles || 0;
    if (availableVehicles) availableVehicles.textContent = stats.available_vehicles || 0;
    if (totalSales) totalSales.textContent = stats.total_sales || 0;
    
    // Показываем выручку только админам
    const userRole = localStorage.getItem('userRole');
    if (userRole === 'admin' && totalRevenue) {
        totalRevenue.textContent = formatCurrency(stats.total_revenue || 0);
        const revenueItem = totalRevenue.closest('.stat-item');
        if (revenueItem) {
            revenueItem.style.display = 'block';
        }
    }
    
    if (totalCustomers) totalCustomers.textContent = stats.total_customers || 0;
}

function formatDateTime(date) {
    return new Intl.DateTimeFormat('ru-RU', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    }).format(new Date(date));
}

// Debounce функция для поиска
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
}

// Копирование в буфер обмена
function copyToClipboard(text) {
    navigator.clipboard.writeText(text).then(() => {
        showNotification('Скопировано в буфер обмена', 'success');
    }).catch(() => {
        showNotification('Ошибка копирования', 'error');
    });
}

// Подтверждение действия
function confirmAction(message, callback) {
    if (confirm(message)) {
        callback();
    }
}

// Экспорт функций для использования в других модулях
window.app = {
    openModal,
    closeModal,
    showNotification,
    showLoading,
    hideLoading,
    formatNumber,
    formatCurrency,
    formatDate,
    formatDateTime,
    debounce,
    copyToClipboard,
    confirmAction
};