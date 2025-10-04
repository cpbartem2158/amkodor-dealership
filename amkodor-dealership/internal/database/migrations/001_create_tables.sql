-- Автоматизированная система автосалона Амкодор (PostgreSQL)
-- База данных состоит из 21 таблицы

-- 1. Таблица категорий техники
CREATE TABLE vehicle_categories (
                                    category_id SERIAL PRIMARY KEY,
                                    category_name VARCHAR(100) NOT NULL UNIQUE,
                                    description TEXT,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    is_active BOOLEAN DEFAULT TRUE
);

-- 2. Таблица типов техники
CREATE TABLE vehicle_types (
                               type_id SERIAL PRIMARY KEY,
                               category_id INTEGER NOT NULL,
                               type_name VARCHAR(100) NOT NULL,
                               description TEXT,
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               FOREIGN KEY (category_id) REFERENCES vehicle_categories(category_id) ON DELETE CASCADE
);

-- 3. Таблица производителей
CREATE TABLE manufacturers (
                               manufacturer_id SERIAL PRIMARY KEY,
                               manufacturer_name VARCHAR(200) NOT NULL,
                               country VARCHAR(100) DEFAULT 'Беларусь',
                               website VARCHAR(200),
                               contact_phone VARCHAR(50),
                               created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 4. Таблица моделей техники
CREATE TABLE vehicle_models (
                                model_id SERIAL PRIMARY KEY,
                                model_name VARCHAR(100) NOT NULL,
                                type_id INTEGER NOT NULL,
                                manufacturer_id INTEGER NOT NULL,
                                description TEXT,
                                specifications JSONB,
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                FOREIGN KEY (type_id) REFERENCES vehicle_types(type_id) ON DELETE CASCADE,
                                FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(manufacturer_id) ON DELETE CASCADE
);

-- 5. Таблица складов/филиалов
CREATE TABLE warehouses (
                            warehouse_id SERIAL PRIMARY KEY,
                            warehouse_name VARCHAR(200) NOT NULL,
                            address TEXT NOT NULL,
                            city VARCHAR(100) NOT NULL,
                            region VARCHAR(100),
                            phone VARCHAR(50),
                            manager_name VARCHAR(200),
                            capacity INTEGER DEFAULT 0,
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                            is_active BOOLEAN DEFAULT TRUE
);

-- 6. Таблица единиц техники
CREATE TABLE vehicles (
                          vehicle_id SERIAL PRIMARY KEY,
                          model_id INTEGER NOT NULL,
                          warehouse_id INTEGER NOT NULL,
                          vin VARCHAR(50) UNIQUE,
                          serial_number VARCHAR(100) NOT NULL,
                          manufacture_year INTEGER NOT NULL,
                          color VARCHAR(50),
                          price DECIMAL(18, 2) NOT NULL CHECK (price >= 0),
                          discount DECIMAL(5, 2) DEFAULT 0 CHECK (discount >= 0 AND discount <= 100),
                          status VARCHAR(50) DEFAULT 'В наличии' CHECK (status IN ('В наличии', 'Продано', 'Зарезервировано', 'В ремонте', 'В пути')),
                          arrival_date DATE DEFAULT CURRENT_DATE,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (model_id) REFERENCES vehicle_models(model_id) ON DELETE CASCADE,
                          FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE CASCADE
);

-- 7. Таблица должностей
CREATE TABLE positions (
                           position_id SERIAL PRIMARY KEY,
                           position_name VARCHAR(100) NOT NULL UNIQUE,
                           base_salary DECIMAL(18, 2) DEFAULT 0 CHECK (base_salary >= 0),
                           description TEXT,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 8. Таблица сотрудников
CREATE TABLE employees (
                           employee_id SERIAL PRIMARY KEY,
                           first_name VARCHAR(100) NOT NULL,
                           last_name VARCHAR(100) NOT NULL,
                           middle_name VARCHAR(100),
                           position_id INTEGER NOT NULL,
                           warehouse_id INTEGER,
                           email VARCHAR(200) UNIQUE,
                           phone VARCHAR(50),
                           password_hash VARCHAR(255) NOT NULL,
                           hire_date DATE DEFAULT CURRENT_DATE,
                           salary DECIMAL(18, 2) CHECK (salary >= 0),
                           is_active BOOLEAN DEFAULT TRUE,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           FOREIGN KEY (position_id) REFERENCES positions(position_id) ON DELETE RESTRICT,
                           FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE SET NULL
);

-- 9. Таблица клиентов
CREATE TABLE customers (
                           customer_id SERIAL PRIMARY KEY,
                           first_name VARCHAR(100) NOT NULL,
                           last_name VARCHAR(100) NOT NULL,
                           middle_name VARCHAR(100),
                           phone VARCHAR(50) NOT NULL,
                           email VARCHAR(200),
                           passport_number VARCHAR(50),
                           address TEXT,
                           date_of_birth DATE,
                           discount_percent DECIMAL(5, 2) DEFAULT 0 CHECK (discount_percent >= 0 AND discount_percent <= 100),
                           is_vip BOOLEAN DEFAULT FALSE,
                           created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                           updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 10. Таблица корпоративных клиентов
CREATE TABLE corporate_clients (
                                   corporate_client_id SERIAL PRIMARY KEY,
                                   company_name VARCHAR(300) NOT NULL,
                                   tax_id VARCHAR(50) UNIQUE NOT NULL,
                                   legal_address TEXT NOT NULL,
                                   contact_person VARCHAR(200),
                                   phone VARCHAR(50) NOT NULL,
                                   email VARCHAR(200),
                                   bank_account VARCHAR(50),
                                   bank_name VARCHAR(200),
                                   discount_percent DECIMAL(5, 2) DEFAULT 0 CHECK (discount_percent >= 0 AND discount_percent <= 100),
                                   contract_number VARCHAR(50),
                                   contract_date DATE,
                                   created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                   updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 11. Таблица продаж
CREATE TABLE sales (
                       sale_id SERIAL PRIMARY KEY,
                       vehicle_id INTEGER NOT NULL,
                       customer_id INTEGER,
                       corporate_client_id INTEGER,
                       employee_id INTEGER NOT NULL,
                       sale_date DATE DEFAULT CURRENT_DATE,
                       base_price DECIMAL(18, 2) NOT NULL CHECK (base_price >= 0),
                       discount_amount DECIMAL(18, 2) DEFAULT 0 CHECK (discount_amount >= 0),
                       final_price DECIMAL(18, 2) NOT NULL CHECK (final_price >= 0),
                       payment_type VARCHAR(50) DEFAULT 'Наличные' CHECK (payment_type IN ('Наличные', 'Безналичный', 'Кредит', 'Лизинг', 'Рассрочка')),
                       status VARCHAR(50) DEFAULT 'Завершена' CHECK (status IN ('Завершена', 'В процессе', 'Отменена')),
                       contract_number VARCHAR(50) UNIQUE,
                       notes TEXT,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE RESTRICT,
                       FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE SET NULL,
                       FOREIGN KEY (corporate_client_id) REFERENCES corporate_clients(corporate_client_id) ON DELETE SET NULL,
                       FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE RESTRICT,
                       CHECK ((customer_id IS NOT NULL AND corporate_client_id IS NULL) OR (customer_id IS NULL AND corporate_client_id IS NOT NULL))
);

-- 12. Таблица тест-драйвов
CREATE TABLE test_drives (
                             test_drive_id SERIAL PRIMARY KEY,
                             vehicle_id INTEGER NOT NULL,
                             customer_id INTEGER,
                             corporate_client_id INTEGER,
                             employee_id INTEGER NOT NULL,
                             scheduled_date TIMESTAMP NOT NULL,
                             duration INTEGER DEFAULT 60,
                             status VARCHAR(50) DEFAULT 'Запланирован' CHECK (status IN ('Запланирован', 'Завершен', 'Отменен', 'Не явился')),
                             feedback_rating INTEGER CHECK (feedback_rating >= 1 AND feedback_rating <= 5),
                             feedback_comment TEXT,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE CASCADE,
                             FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE SET NULL,
                             FOREIGN KEY (corporate_client_id) REFERENCES corporate_clients(corporate_client_id) ON DELETE SET NULL,
                             FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE RESTRICT,
                             CHECK ((customer_id IS NOT NULL AND corporate_client_id IS NULL) OR (customer_id IS NULL AND corporate_client_id IS NOT NULL)),
                             CHECK (scheduled_date > CURRENT_TIMESTAMP)
);

-- 13. Таблица сервисного обслуживания
CREATE TABLE service_orders (
                                service_order_id SERIAL PRIMARY KEY,
                                vehicle_id INTEGER NOT NULL,
                                customer_id INTEGER,
                                corporate_client_id INTEGER,
                                employee_id INTEGER NOT NULL,
                                order_date DATE DEFAULT CURRENT_DATE,
                                completion_date DATE,
                                service_type VARCHAR(100) NOT NULL,
                                description TEXT,
                                cost DECIMAL(18, 2) DEFAULT 0 CHECK (cost >= 0),
                                status VARCHAR(50) DEFAULT 'В работе' CHECK (status IN ('В работе', 'Завершен', 'Приостановлен', 'Отменен')),
                                created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                FOREIGN KEY (vehicle_id) REFERENCES vehicles(vehicle_id) ON DELETE RESTRICT,
                                FOREIGN KEY (customer_id) REFERENCES customers(customer_id) ON DELETE SET NULL,
                                FOREIGN KEY (corporate_client_id) REFERENCES corporate_clients(corporate_client_id) ON DELETE SET NULL,
                                FOREIGN KEY (employee_id) REFERENCES employees(employee_id) ON DELETE RESTRICT,
                                CHECK ((customer_id IS NOT NULL AND corporate_client_id IS NULL) OR (customer_id IS NULL AND corporate_client_id IS NOT NULL)),
                                CHECK (completion_date IS NULL OR completion_date >= order_date)
);

-- 14. Таблица запасных частей
CREATE TABLE spare_parts (
                             spare_part_id SERIAL PRIMARY KEY,
                             part_number VARCHAR(100) NOT NULL UNIQUE,
                             part_name VARCHAR(200) NOT NULL,
                             model_id INTEGER,
                             price DECIMAL(18, 2) NOT NULL CHECK (price >= 0),
                             quantity_in_stock INTEGER DEFAULT 0 CHECK (quantity_in_stock >= 0),
                             min_quantity INTEGER DEFAULT 0,
                             warehouse_id INTEGER NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                             FOREIGN KEY (model_id) REFERENCES vehicle_models(model_id) ON DELETE SET NULL,
                             FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE RESTRICT
);

-- 15. Таблица использования запчастей
CREATE TABLE service_order_parts (
                                     service_order_part_id SERIAL PRIMARY KEY,
                                     service_order_id INTEGER NOT NULL,
                                     spare_part_id INTEGER NOT NULL,
                                     quantity INTEGER NOT NULL CHECK (quantity > 0),
                                     unit_price DECIMAL(18, 2) NOT NULL CHECK (unit_price >= 0),
                                     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                     FOREIGN KEY (service_order_id) REFERENCES service_orders(service_order_id) ON DELETE CASCADE,
                                     FOREIGN KEY (spare_part_id) REFERENCES spare_parts(spare_part_id) ON DELETE RESTRICT
);

-- 16. Таблица поставок техники
CREATE TABLE supplies (
                          supply_id SERIAL PRIMARY KEY,
                          manufacturer_id INTEGER NOT NULL,
                          warehouse_id INTEGER NOT NULL,
                          supply_date DATE DEFAULT CURRENT_DATE,
                          expected_arrival_date DATE,
                          actual_arrival_date DATE,
                          total_cost DECIMAL(18, 2) DEFAULT 0 CHECK (total_cost >= 0),
                          status VARCHAR(50) DEFAULT 'В пути' CHECK (status IN ('Заказано', 'В пути', 'Получено', 'Отменено')),
                          invoice_number VARCHAR(50),
                          notes TEXT,
                          created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          FOREIGN KEY (manufacturer_id) REFERENCES manufacturers(manufacturer_id) ON DELETE RESTRICT,
                          FOREIGN KEY (warehouse_id) REFERENCES warehouses(warehouse_id) ON DELETE RESTRICT,
                          CHECK (expected_arrival_date >= supply_date),
                          CHECK (actual_arrival_date IS NULL OR actual_arrival_date >= supply_date)
);

-- 17. Таблица элементов поставки
CREATE TABLE supply_items (
                              supply_item_id SERIAL PRIMARY KEY,
                              supply_id INTEGER NOT NULL,
                              model_id INTEGER NOT NULL,
                              quantity INTEGER NOT NULL CHECK (quantity > 0),
                              unit_price DECIMAL(18, 2) NOT NULL CHECK (unit_price >= 0),
                              created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                              FOREIGN KEY (supply_id) REFERENCES supplies(supply_id) ON DELETE CASCADE,
                              FOREIGN KEY (model_id) REFERENCES vehicle_models(model_id) ON DELETE RESTRICT
);

-- 18. History таблица для продаж
CREATE TABLE sales_history (
                               history_id SERIAL PRIMARY KEY,
                               sale_id INTEGER NOT NULL,
                               operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('INSERT', 'UPDATE', 'DELETE')),
                               operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                               old_value JSONB,
                               new_value JSONB,
                               username VARCHAR(200) DEFAULT CURRENT_USER,
                               hostname VARCHAR(200),
                               application_name VARCHAR(200)
);

-- 19. History таблица для техники
CREATE TABLE vehicles_history (
                                  history_id SERIAL PRIMARY KEY,
                                  vehicle_id INTEGER NOT NULL,
                                  operation_type VARCHAR(20) NOT NULL CHECK (operation_type IN ('INSERT', 'UPDATE', 'DELETE')),
                                  operation_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                  old_value JSONB,
                                  new_value JSONB,
                                  username VARCHAR(200) DEFAULT CURRENT_USER,
                                  hostname VARCHAR(200),
                                  application_name VARCHAR(200)
);

-- 20. Таблица логирования SSIS/ETL
CREATE TABLE ssis_log (
                          log_id SERIAL PRIMARY KEY,
                          package_name VARCHAR(200) NOT NULL,
                          execution_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                          rows_processed INTEGER DEFAULT 0,
                          status VARCHAR(50) DEFAULT 'Success' CHECK (status IN ('Success', 'Failed', 'Warning')),
                          error_message TEXT,
                          username VARCHAR(200) DEFAULT CURRENT_USER,
                          execution_duration INTEGER
);

-- 21. Таблица конфигурации SSIS/ETL
CREATE TABLE ssis_configuration (
                                    config_id SERIAL PRIMARY KEY,
                                    config_name VARCHAR(100) NOT NULL UNIQUE,
                                    config_value TEXT,
                                    description TEXT,
                                    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                                    updated_by VARCHAR(200) DEFAULT CURRENT_USER
);

-- Индексы для оптимизации
CREATE INDEX idx_vehicles_status ON vehicles(status);
CREATE INDEX idx_vehicles_model_id ON vehicles(model_id);
CREATE INDEX idx_sales_sale_date ON sales(sale_date);
CREATE INDEX idx_sales_employee_id ON sales(employee_id);
CREATE INDEX idx_customers_phone ON customers(phone);
CREATE INDEX idx_customers_email ON customers(email);
CREATE INDEX idx_employees_position_id ON employees(position_id);
CREATE INDEX idx_employees_email ON employees(email);
CREATE INDEX idx_test_drives_scheduled_date ON test_drives(scheduled_date);
CREATE INDEX idx_service_orders_status ON service_orders(status);
CREATE INDEX idx_spare_parts_part_number ON spare_parts(part_number);

-- Функция обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
RETURN NEW;
END;
$$ language 'plpgsql';

-- Триггеры для автоматического обновления updated_at
CREATE TRIGGER update_vehicles_updated_at BEFORE UPDATE ON vehicles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_customers_updated_at BEFORE UPDATE ON customers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_corporate_clients_updated_at BEFORE UPDATE ON corporate_clients
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sales_updated_at BEFORE UPDATE ON sales
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();