-- Хранимые процедуры для автосалона Амкодор (PostgreSQL)

-- 1. Процедура поиска техники по различным критериям
CREATE OR REPLACE FUNCTION sp_search_vehicles(
    p_model_name VARCHAR(100) DEFAULT NULL,
    p_category_name VARCHAR(100) DEFAULT NULL,
    p_type_name VARCHAR(100) DEFAULT NULL,
    p_manufacturer_name VARCHAR(200) DEFAULT NULL,
    p_min_price DECIMAL(18, 2) DEFAULT NULL,
    p_max_price DECIMAL(18, 2) DEFAULT NULL,
    p_min_year INTEGER DEFAULT NULL,
    p_max_year INTEGER DEFAULT NULL,
    p_status VARCHAR(50) DEFAULT NULL,
    p_warehouse_id INTEGER DEFAULT NULL,
    p_city VARCHAR(100) DEFAULT NULL
)
    RETURNS TABLE (
                      vehicle_id INTEGER,
                      vin VARCHAR(50),
                      serial_number VARCHAR(100),
                      model_name VARCHAR(100),
                      type_name VARCHAR(100),
                      category_name VARCHAR(100),
                      manufacturer_name VARCHAR(200),
                      manufacture_year INTEGER,
                      color VARCHAR(50),
                      price DECIMAL(18, 2),
                      discount DECIMAL(5, 2),
                      final_price DECIMAL(18, 2),
                      status VARCHAR(50),
                      warehouse_name VARCHAR(200),
                      city VARCHAR(100),
                      warehouse_phone VARCHAR(50),
                      description TEXT,
                      specifications JSONB
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            v.vehicle_id,
            v.vin,
            v.serial_number,
            vm.model_name,
            vt.type_name,
            vc.category_name,
            m.manufacturer_name,
            v.manufacture_year,
            v.color,
            v.price,
            v.discount,
            fn_calculate_final_price(v.price, v.discount) AS final_price,
            v.status,
            w.warehouse_name,
            w.city,
            w.phone AS warehouse_phone,
            vm.description,
            vm.specifications
        FROM vehicles v
                 INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
                 INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
                 INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
                 INNER JOIN manufacturers m ON vm.manufacturer_id = m.manufacturer_id
                 INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
        WHERE
            (p_model_name IS NULL OR vm.model_name ILIKE '%' || p_model_name || '%')
          AND (p_category_name IS NULL OR vc.category_name ILIKE '%' || p_category_name || '%')
          AND (p_type_name IS NULL OR vt.type_name ILIKE '%' || p_type_name || '%')
          AND (p_manufacturer_name IS NULL OR m.manufacturer_name ILIKE '%' || p_manufacturer_name || '%')
          AND (p_min_price IS NULL OR v.price >= p_min_price)
          AND (p_max_price IS NULL OR v.price <= p_max_price)
          AND (p_min_year IS NULL OR v.manufacture_year >= p_min_year)
          AND (p_max_year IS NULL OR v.manufacture_year <= p_max_year)
          AND (p_status IS NULL OR v.status = p_status)
          AND (p_warehouse_id IS NULL OR v.warehouse_id = p_warehouse_id)
          AND (p_city IS NULL OR w.city ILIKE '%' || p_city || '%')
        ORDER BY v.created_at DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- 2. Процедура поиска клиентов
CREATE OR REPLACE FUNCTION sp_search_customers(
    p_search_term VARCHAR(200) DEFAULT NULL,
    p_phone VARCHAR(50) DEFAULT NULL,
    p_email VARCHAR(200) DEFAULT NULL,
    p_is_vip BOOLEAN DEFAULT NULL,
    p_min_discount DECIMAL(5, 2) DEFAULT NULL
)
    RETURNS TABLE (
                      customer_id INTEGER,
                      full_name VARCHAR(300),
                      phone VARCHAR(50),
                      email VARCHAR(200),
                      address TEXT,
                      discount_percent DECIMAL(5, 2),
                      is_vip BOOLEAN,
                      customer_level VARCHAR(50),
                      total_purchases BIGINT,
                      total_spent DECIMAL(18, 2),
                      created_at TIMESTAMP
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            c.customer_id,
            c.last_name || ' ' || c.first_name || COALESCE(' ' || c.middle_name, '') AS full_name,
            c.phone,
            c.email,
            c.address,
            c.discount_percent,
            c.is_vip,
            fn_get_customer_level(c.customer_id) AS customer_level,
            (SELECT COUNT(*) FROM sales s WHERE s.customer_id = c.customer_id AND s.status = 'Завершена')::BIGINT AS total_purchases,
            (SELECT COALESCE(SUM(s.final_price), 0) FROM sales s WHERE s.customer_id = c.customer_id AND s.status = 'Завершена') AS total_spent,
            c.created_at
        FROM customers c
        WHERE
            (p_search_term IS NULL OR
             c.last_name ILIKE '%' || p_search_term || '%' OR
             c.first_name ILIKE '%' || p_search_term || '%' OR
             c.middle_name ILIKE '%' || p_search_term || '%')
          AND (p_phone IS NULL OR c.phone ILIKE '%' || p_phone || '%')
          AND (p_email IS NULL OR c.email ILIKE '%' || p_email || '%')
          AND (p_is_vip IS NULL OR c.is_vip = p_is_vip)
          AND (p_min_discount IS NULL OR c.discount_percent >= p_min_discount)
        ORDER BY c.last_name, c.first_name;
END;
$$ LANGUAGE plpgsql STABLE;

-- 3. Процедура формирования отчета о продажах (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_generate_sales_report(
    p_start_date DATE,
    p_end_date DATE,
    p_employee_id INTEGER DEFAULT NULL,
    p_warehouse_id INTEGER DEFAULT NULL,
    p_category_id INTEGER DEFAULT NULL
)
    RETURNS TABLE (
                      sale_id INTEGER,
                      contract_number VARCHAR(50),
                      sale_date DATE,
                      model_name VARCHAR(100),
                      type_name VARCHAR(100),
                      category_name VARCHAR(100),
                      client_name VARCHAR(300),
                      manager_name VARCHAR(200),
                      warehouse_name VARCHAR(200),
                      city VARCHAR(100),
                      base_price DECIMAL(18, 2),
                      discount_amount DECIMAL(18, 2),
                      final_price DECIMAL(18, 2),
                      payment_type VARCHAR(50),
                      status VARCHAR(50)
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            s.sale_id,
            s.contract_number,
            s.sale_date,
            vm.model_name,
            vt.type_name,
            vc.category_name,
            fn_get_client_full_name(s.customer_id, s.corporate_client_id) AS client_name,
            e.last_name || ' ' || e.first_name AS manager_name,
            w.warehouse_name,
            w.city,
            s.base_price,
            s.discount_amount,
            s.final_price,
            s.payment_type,
            s.status
        FROM sales s
                 INNER JOIN vehicles v ON s.vehicle_id = v.vehicle_id
                 INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
                 INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
                 INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
                 INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
                 INNER JOIN employees e ON s.employee_id = e.employee_id
        WHERE
            s.sale_date BETWEEN p_start_date AND p_end_date
          AND (p_employee_id IS NULL OR s.employee_id = p_employee_id)
          AND (p_warehouse_id IS NULL OR v.warehouse_id = p_warehouse_id)
          AND (p_category_id IS NULL OR vc.category_id = p_category_id)
          AND s.status = 'Завершена'
        ORDER BY s.sale_date DESC;
END;
$$ LANGUAGE plpgsql STABLE;

-- 4. Процедура формирования отчета по инвентарю (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_generate_inventory_report(
    p_warehouse_id INTEGER DEFAULT NULL,
    p_category_id INTEGER DEFAULT NULL,
    p_status VARCHAR(50) DEFAULT NULL
)
    RETURNS TABLE (
                      warehouse_name VARCHAR(200),
                      city VARCHAR(100),
                      category_name VARCHAR(100),
                      type_name VARCHAR(100),
                      model_name VARCHAR(100),
                      quantity BIGINT,
                      status VARCHAR(50),
                      average_price DECIMAL(18, 2),
                      total_value DECIMAL(18, 2),
                      oldest_arrival_date DATE,
                      newest_arrival_date DATE
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            w.warehouse_name,
            w.city,
            vc.category_name,
            vt.type_name,
            vm.model_name,
            COUNT(v.vehicle_id) AS quantity,
            v.status,
            AVG(v.price) AS average_price,
            SUM(fn_calculate_final_price(v.price, v.discount)) AS total_value,
            MIN(v.arrival_date) AS oldest_arrival_date,
            MAX(v.arrival_date) AS newest_arrival_date
        FROM vehicles v
                 INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
                 INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
                 INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
                 INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
        WHERE
            (p_warehouse_id IS NULL OR v.warehouse_id = p_warehouse_id)
          AND (p_category_id IS NULL OR vc.category_id = p_category_id)
          AND (p_status IS NULL OR v.status = p_status)
        GROUP BY
            w.warehouse_name,
            w.city,
            vc.category_name,
            vt.type_name,
            vm.model_name,
            v.status
        ORDER BY
            w.warehouse_name,
            vc.category_name,
            vm.model_name;
END;
$$ LANGUAGE plpgsql STABLE;

-- 5. Процедура создания продажи (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_create_sale(
    p_vehicle_id INTEGER,
    p_customer_id INTEGER DEFAULT NULL,
    p_corporate_client_id INTEGER DEFAULT NULL,
    p_employee_id INTEGER DEFAULT NULL,
    p_payment_type VARCHAR(50) DEFAULT 'Наличные',
    p_additional_discount DECIMAL(5, 2) DEFAULT 0,
    p_contract_number VARCHAR(50) DEFAULT NULL,
    p_notes TEXT DEFAULT NULL
)
    RETURNS INTEGER AS $$
DECLARE
    v_sale_id INTEGER;
    v_base_price DECIMAL(18, 2);
    v_vehicle_discount DECIMAL(5, 2);
    v_client_discount DECIMAL(5, 2) := 0;
    v_total_discount DECIMAL(5, 2);
    v_discount_amount DECIMAL(18, 2);
    v_final_price DECIMAL(18, 2);
BEGIN
    -- Проверка что техника доступна
    IF NOT fn_is_vehicle_available(p_vehicle_id) THEN
        RAISE EXCEPTION 'Техника недоступна для продажи';
    END IF;

    -- Получение базовой цены и скидки техники
    SELECT price, discount
    INTO v_base_price, v_vehicle_discount
    FROM vehicles
    WHERE vehicle_id = p_vehicle_id;

    -- Получение скидки клиента
    IF p_customer_id IS NOT NULL THEN
        SELECT discount_percent INTO v_client_discount
        FROM customers WHERE customer_id = p_customer_id;
    ELSIF p_corporate_client_id IS NOT NULL THEN
        SELECT discount_percent INTO v_client_discount
        FROM corporate_clients WHERE corporate_client_id = p_corporate_client_id;
    END IF;

    -- Расчет общей скидки
    v_total_discount := v_vehicle_discount + v_client_discount + p_additional_discount;
    IF v_total_discount > 100 THEN v_total_discount := 100; END IF;

    v_discount_amount := fn_calculate_discount_amount(v_base_price, v_total_discount);
    v_final_price := fn_calculate_final_price(v_base_price, v_total_discount);

    -- Создание продажи
    INSERT INTO sales (
        vehicle_id, customer_id, corporate_client_id, employee_id,
        base_price, discount_amount, final_price, payment_type,
        contract_number, notes
    ) VALUES (
                 p_vehicle_id, p_customer_id, p_corporate_client_id, p_employee_id,
                 v_base_price, v_discount_amount, v_final_price, p_payment_type,
                 p_contract_number, p_notes
             ) RETURNING sale_id INTO v_sale_id;

    -- Обновление статуса техники
    UPDATE vehicles
    SET status = 'Продано', updated_at = CURRENT_TIMESTAMP
    WHERE vehicle_id = p_vehicle_id;

    RETURN v_sale_id;
END;
$$ LANGUAGE plpgsql;

-- 6. Процедура создания тест-драйва (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_create_test_drive(
    p_vehicle_id INTEGER,
    p_customer_id INTEGER DEFAULT NULL,
    p_corporate_client_id INTEGER DEFAULT NULL,
    p_employee_id INTEGER DEFAULT NULL,
    p_scheduled_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP + INTERVAL '1 day',
    p_duration INTEGER DEFAULT 60
)
    RETURNS INTEGER AS $$
DECLARE
    v_test_drive_id INTEGER;
BEGIN
    -- Проверка что дата в будущем
    IF p_scheduled_date <= CURRENT_TIMESTAMP THEN
        RAISE EXCEPTION 'Дата тест-драйва должна быть в будущем';
    END IF;

    -- Проверка доступности техники
    IF NOT fn_is_vehicle_available(p_vehicle_id) THEN
        RAISE EXCEPTION 'Техника недоступна для тест-драйва';
    END IF;

    -- Создание тест-драйва
    INSERT INTO test_drives (
        vehicle_id, customer_id, corporate_client_id, employee_id,
        scheduled_date, duration
    ) VALUES (
                 p_vehicle_id, p_customer_id, p_corporate_client_id, p_employee_id,
                 p_scheduled_date, p_duration
             ) RETURNING test_drive_id INTO v_test_drive_id;

    RETURN v_test_drive_id;
END;
$$ LANGUAGE plpgsql;

-- 7. Процедура создания сервисного заказа (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_create_service_order(
    p_vehicle_id INTEGER,
    p_employee_id INTEGER,           -- обязательные параметры первыми
    p_service_type VARCHAR(100),     -- обязательные параметры первыми
    p_customer_id INTEGER DEFAULT NULL,
    p_corporate_client_id INTEGER DEFAULT NULL,
    p_description TEXT DEFAULT NULL,
    p_cost DECIMAL(18, 2) DEFAULT 0
)
    RETURNS INTEGER AS $$
DECLARE
    v_service_order_id INTEGER;
BEGIN
    INSERT INTO service_orders (
        vehicle_id, customer_id, corporate_client_id, employee_id,
        service_type, description, cost
    ) VALUES (
                 p_vehicle_id, p_customer_id, p_corporate_client_id, p_employee_id,
                 p_service_type, p_description, p_cost
             ) RETURNING service_order_id INTO v_service_order_id;

    RETURN v_service_order_id;
END;
$$ LANGUAGE plpgsql;

-- 8. Процедура добавления запчастей в сервисный заказ (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_add_spare_parts_to_service(
    p_service_order_id INTEGER,
    p_spare_part_id INTEGER,
    p_quantity INTEGER
)
    RETURNS INTEGER AS $$
DECLARE
    v_service_order_part_id INTEGER;
    v_unit_price DECIMAL(18, 2);
    v_available_quantity INTEGER;
BEGIN
    -- Проверка доступности запчастей
    SELECT price, quantity_in_stock
    INTO v_unit_price, v_available_quantity
    FROM spare_parts
    WHERE spare_part_id = p_spare_part_id;

    IF v_available_quantity < p_quantity THEN
        RAISE EXCEPTION 'Недостаточно запчастей на складе';
    END IF;

    -- Добавление запчастей в заказ
    INSERT INTO service_order_parts (
        service_order_id, spare_part_id, quantity, unit_price
    ) VALUES (
                 p_service_order_id, p_spare_part_id, p_quantity, v_unit_price
             ) RETURNING service_order_part_id INTO v_service_order_part_id;

    -- Обновление остатков
    UPDATE spare_parts
    SET quantity_in_stock = quantity_in_stock - p_quantity
    WHERE spare_part_id = p_spare_part_id;

    RETURN v_service_order_part_id;
END;
$$ LANGUAGE plpgsql;

-- 9. Процедура завершения сервисного заказа (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_complete_service_order(
    p_service_order_id INTEGER
)
    RETURNS VOID AS $$
DECLARE
    v_parts_cost DECIMAL(18, 2);
    v_service_cost DECIMAL(18, 2);
    v_total_cost DECIMAL(18, 2);
BEGIN
    -- Расчет стоимости запчастей
    SELECT COALESCE(SUM(quantity * unit_price), 0)
    INTO v_parts_cost
    FROM service_order_parts
    WHERE service_order_id = p_service_order_id;

    -- Получение стоимости услуги
    SELECT cost INTO v_service_cost
    FROM service_orders
    WHERE service_order_id = p_service_order_id;

    v_total_cost := v_service_cost + v_parts_cost;

    -- Обновление заказа
    UPDATE service_orders
    SET
        cost = v_total_cost,
        status = 'Завершен',
        completion_date = CURRENT_DATE
    WHERE service_order_id = p_service_order_id;
END;
$$ LANGUAGE plpgsql;

-- 10. Процедура получения статистики продаж (ИСПРАВЛЕНА)
CREATE OR REPLACE FUNCTION sp_get_sales_statistics(
    p_start_date DATE,
    p_end_date DATE
)
    RETURNS TABLE (
                      total_sales BIGINT,
                      total_revenue DECIMAL(18, 2),
                      average_sale_price DECIMAL(18, 2),
                      total_discounts DECIMAL(18, 2),
                      best_category VARCHAR(100),
                      best_manager VARCHAR(200)
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH stats AS (
            SELECT
                COUNT(*) as sales_count,
                SUM(s.final_price) as revenue,
                AVG(s.final_price) as avg_price,
                SUM(s.discount_amount) as discounts
            FROM sales s
            WHERE s.sale_date BETWEEN p_start_date AND p_end_date
              AND s.status = 'Завершена'
        ),
             best_cat AS (
                 SELECT vc.category_name
                 FROM sales s
                          INNER JOIN vehicles v ON s.vehicle_id = v.vehicle_id
                          INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
                          INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
                          INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
                 WHERE s.sale_date BETWEEN p_start_date AND p_end_date
                   AND s.status = 'Завершена'
                 GROUP BY vc.category_name
                 ORDER BY COUNT(*) DESC
                 LIMIT 1
             ),
             best_mgr AS (
                 SELECT e.last_name || ' ' || e.first_name as manager_name
                 FROM sales s
                          INNER JOIN employees e ON s.employee_id = e.employee_id
                 WHERE s.sale_date BETWEEN p_start_date AND p_end_date
                   AND s.status = 'Завершена'
                 GROUP BY e.employee_id, e.last_name, e.first_name
                 ORDER BY SUM(s.final_price) DESC
                 LIMIT 1
             )
        SELECT
            st.sales_count::BIGINT,
            st.revenue,
            st.avg_price,
            st.discounts,
            bc.category_name,
            bm.manager_name
        FROM stats st
                 CROSS JOIN best_cat bc
                 CROSS JOIN best_mgr bm;
END;
$$ LANGUAGE plpgsql STABLE;