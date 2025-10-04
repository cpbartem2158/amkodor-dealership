-- Функции для автосалона Амкодор (PostgreSQL)

-- 1. Функция расчета финальной цены с учетом скидки
CREATE OR REPLACE FUNCTION fn_calculate_final_price(
    p_base_price DECIMAL(18, 2),
    p_discount_percent DECIMAL(5, 2)
)
    RETURNS DECIMAL(18, 2) AS $$
DECLARE
    v_final_price DECIMAL(18, 2);
BEGIN
    -- Проверка границ скидки
    IF p_discount_percent < 0 THEN p_discount_percent := 0; END IF;
    IF p_discount_percent > 100 THEN p_discount_percent := 100; END IF;

    v_final_price := ROUND(p_base_price * (1 - p_discount_percent / 100), 2);

    RETURN v_final_price;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- 2. Функция расчета скидки в денежном выражении
CREATE OR REPLACE FUNCTION fn_calculate_discount_amount(
    p_base_price DECIMAL(18, 2),
    p_discount_percent DECIMAL(5, 2)
)
    RETURNS DECIMAL(18, 2) AS $$
DECLARE
    v_discount_amount DECIMAL(18, 2);
BEGIN
    IF p_discount_percent < 0 THEN p_discount_percent := 0; END IF;
    IF p_discount_percent > 100 THEN p_discount_percent := 100; END IF;

    v_discount_amount := ROUND(p_base_price * (p_discount_percent / 100), 2);

    RETURN v_discount_amount;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- 3. Функция подсчета общей прибыли за период
CREATE OR REPLACE FUNCTION fn_calculate_profit_by_period(
    p_start_date DATE,
    p_end_date DATE
)
    RETURNS DECIMAL(18, 2) AS $$
DECLARE
    v_total_profit DECIMAL(18, 2);
BEGIN
    SELECT COALESCE(SUM(final_price), 0)
    INTO v_total_profit
    FROM sales
    WHERE sale_date BETWEEN p_start_date AND p_end_date
      AND status = 'Завершена';

    RETURN v_total_profit;
END;
$$ LANGUAGE plpgsql STABLE;

-- 4. Функция подсчета продаж по менеджеру
CREATE OR REPLACE FUNCTION fn_count_sales_by_employee(
    p_employee_id INTEGER,
    p_start_date DATE DEFAULT NULL,
    p_end_date DATE DEFAULT NULL
)
    RETURNS INTEGER AS $$
DECLARE
    v_sales_count INTEGER;
    v_start DATE;
    v_end DATE;
BEGIN
    v_start := COALESCE(p_start_date, '1900-01-01');
    v_end := COALESCE(p_end_date, '9999-12-31');

    SELECT COUNT(*)
    INTO v_sales_count
    FROM sales
    WHERE employee_id = p_employee_id
      AND sale_date BETWEEN v_start AND v_end
      AND status = 'Завершена';

    RETURN COALESCE(v_sales_count, 0);
END;
$$ LANGUAGE plpgsql STABLE;

-- 5. Функция расчета комиссии менеджера
CREATE OR REPLACE FUNCTION fn_calculate_manager_commission(
    p_employee_id INTEGER,
    p_commission_percent DECIMAL(5, 2),
    p_start_date DATE,
    p_end_date DATE
)
    RETURNS DECIMAL(18, 2) AS $$
DECLARE
    v_commission DECIMAL(18, 2);
    v_total_sales DECIMAL(18, 2);
BEGIN
    SELECT COALESCE(SUM(final_price), 0)
    INTO v_total_sales
    FROM sales
    WHERE employee_id = p_employee_id
      AND sale_date BETWEEN p_start_date AND p_end_date
      AND status = 'Завершена';

    v_commission := ROUND(v_total_sales * (p_commission_percent / 100), 2);

    RETURN v_commission;
END;
$$ LANGUAGE plpgsql STABLE;

-- 6. Функция подсчета техники на складе
CREATE OR REPLACE FUNCTION fn_count_vehicles_in_warehouse(
    p_warehouse_id INTEGER,
    p_status VARCHAR(50) DEFAULT NULL
)
    RETURNS INTEGER AS $$
DECLARE
    v_vehicle_count INTEGER;
BEGIN
    IF p_status IS NULL THEN
        SELECT COUNT(*)
        INTO v_vehicle_count
        FROM vehicles
        WHERE warehouse_id = p_warehouse_id;
    ELSE
        SELECT COUNT(*)
        INTO v_vehicle_count
        FROM vehicles
        WHERE warehouse_id = p_warehouse_id
          AND status = p_status;
    END IF;

    RETURN COALESCE(v_vehicle_count, 0);
END;
$$ LANGUAGE plpgsql STABLE;

-- 7. Функция расчета стоимости инвентаря на складе
CREATE OR REPLACE FUNCTION fn_calculate_warehouse_inventory_value(
    p_warehouse_id INTEGER
)
    RETURNS DECIMAL(18, 2) AS $$
DECLARE
    v_total_value DECIMAL(18, 2);
BEGIN
    SELECT COALESCE(SUM(ROUND(price * (1 - discount / 100), 2)), 0)
    INTO v_total_value
    FROM vehicles
    WHERE warehouse_id = p_warehouse_id
      AND status IN ('В наличии', 'Зарезервировано');

    RETURN v_total_value;
END;
$$ LANGUAGE plpgsql STABLE;

-- 8. Функция определения уровня клиента
CREATE OR REPLACE FUNCTION fn_get_customer_level(
    p_customer_id INTEGER
)
    RETURNS VARCHAR(50) AS $$
DECLARE
    v_level VARCHAR(50);
    v_total_purchases DECIMAL(18, 2);
    v_purchase_count INTEGER;
BEGIN
    SELECT
        COALESCE(SUM(final_price), 0),
        COUNT(*)
    INTO v_total_purchases, v_purchase_count
    FROM sales
    WHERE customer_id = p_customer_id
      AND status = 'Завершена';

    IF v_purchase_count = 0 THEN
        v_level := 'Новый';
    ELSIF v_total_purchases >= 1000000 OR v_purchase_count >= 5 THEN
        v_level := 'Платиновый';
    ELSIF v_total_purchases >= 500000 OR v_purchase_count >= 3 THEN
        v_level := 'Золотой';
    ELSIF v_total_purchases >= 100000 OR v_purchase_count >= 1 THEN
        v_level := 'Серебряный';
    ELSE
        v_level := 'Бронзовый';
    END IF;

    RETURN v_level;
END;
$$ LANGUAGE plpgsql STABLE;

-- 9. Функция подсчета завершенных сервисных заказов
CREATE OR REPLACE FUNCTION fn_count_completed_service_orders(
    p_vehicle_id INTEGER DEFAULT NULL,
    p_start_date DATE DEFAULT NULL,
    p_end_date DATE DEFAULT NULL
)
    RETURNS INTEGER AS $$
DECLARE
    v_order_count INTEGER;
    v_start DATE;
    v_end DATE;
BEGIN
    v_start := COALESCE(p_start_date, '1900-01-01');
    v_end := COALESCE(p_end_date, '9999-12-31');

    IF p_vehicle_id IS NULL THEN
        SELECT COUNT(*)
        INTO v_order_count
        FROM service_orders
        WHERE status = 'Завершен'
          AND order_date BETWEEN v_start AND v_end;
    ELSE
        SELECT COUNT(*)
        INTO v_order_count
        FROM service_orders
        WHERE vehicle_id = p_vehicle_id
          AND status = 'Завершен'
          AND order_date BETWEEN v_start AND v_end;
    END IF;

    RETURN COALESCE(v_order_count, 0);
END;
$$ LANGUAGE plpgsql STABLE;

-- 10. Функция расчета среднего рейтинга тест-драйвов
CREATE OR REPLACE FUNCTION fn_get_average_test_drive_rating(
    p_model_id INTEGER
)
    RETURNS DECIMAL(3, 2) AS $$
DECLARE
    v_avg_rating DECIMAL(3, 2);
BEGIN
    SELECT AVG(feedback_rating)
    INTO v_avg_rating
    FROM test_drives td
             INNER JOIN vehicles v ON td.vehicle_id = v.vehicle_id
    WHERE v.model_id = p_model_id
      AND td.status = 'Завершен'
      AND td.feedback_rating IS NOT NULL;

    RETURN COALESCE(v_avg_rating, 0);
END;
$$ LANGUAGE plpgsql STABLE;

-- 11. Функция проверки доступности техники
CREATE OR REPLACE FUNCTION fn_is_vehicle_available(
    p_vehicle_id INTEGER
)
    RETURNS BOOLEAN AS $$
DECLARE
    v_is_available BOOLEAN;
BEGIN
    SELECT CASE
               WHEN status = 'В наличии' THEN TRUE
               ELSE FALSE
               END
    INTO v_is_available
    FROM vehicles
    WHERE vehicle_id = p_vehicle_id;

    RETURN COALESCE(v_is_available, FALSE);
END;
$$ LANGUAGE plpgsql STABLE;

-- 12. Функция получения полного имени клиента
CREATE OR REPLACE FUNCTION fn_get_client_full_name(
    p_customer_id INTEGER DEFAULT NULL,
    p_corporate_client_id INTEGER DEFAULT NULL
)
    RETURNS VARCHAR(300) AS $$
DECLARE
    v_full_name VARCHAR(300);
BEGIN
    IF p_customer_id IS NOT NULL THEN
        SELECT last_name || ' ' || first_name || COALESCE(' ' || middle_name, '')
        INTO v_full_name
        FROM customers
        WHERE customer_id = p_customer_id;
    ELSIF p_corporate_client_id IS NOT NULL THEN
        SELECT company_name
        INTO v_full_name
        FROM corporate_clients
        WHERE corporate_client_id = p_corporate_client_id;
    END IF;

    RETURN COALESCE(v_full_name, 'Не указано');
END;
$$ LANGUAGE plpgsql STABLE;

-- 13. Функция расчета возраста техники
CREATE OR REPLACE FUNCTION fn_get_vehicle_age(
    p_manufacture_year INTEGER
)
    RETURNS INTEGER AS $$
BEGIN
    RETURN EXTRACT(YEAR FROM CURRENT_DATE) - p_manufacture_year;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

-- 14. Функция подсчета дней на складе
CREATE OR REPLACE FUNCTION fn_get_days_in_warehouse(
    p_vehicle_id INTEGER
)
    RETURNS INTEGER AS $$
DECLARE
    v_days INTEGER;
BEGIN
    SELECT CURRENT_DATE - arrival_date
    INTO v_days
    FROM vehicles
    WHERE vehicle_id = p_vehicle_id;

    RETURN COALESCE(v_days, 0);
END;
$$ LANGUAGE plpgsql STABLE;

-- 15. Функция определения статуса запаса запчастей
CREATE OR REPLACE FUNCTION fn_get_spare_part_stock_status(
    p_spare_part_id INTEGER
)
    RETURNS VARCHAR(50) AS $$
DECLARE
    v_status VARCHAR(50);
    v_quantity INTEGER;
    v_min_quantity INTEGER;
BEGIN
    SELECT
        quantity_in_stock,
        min_quantity
    INTO v_quantity, v_min_quantity
    FROM spare_parts
    WHERE spare_part_id = p_spare_part_id;

    IF v_quantity = 0 THEN
        v_status := 'Нет в наличии';
    ELSIF v_quantity <= v_min_quantity THEN
        v_status := 'Критический остаток';
    ELSIF v_quantity <= v_min_quantity * 2 THEN
        v_status := 'Низкий остаток';
    ELSE
        v_status := 'В наличии';
    END IF;

    RETURN v_status;
END;
$$ LANGUAGE plpgsql STABLE;