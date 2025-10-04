-- Триггеры для автосалона Амкодор (PostgreSQL)

-- Функция для логирования изменений в таблице sales
CREATE OR REPLACE FUNCTION log_sales_changes()
    RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO sales_history (
            sale_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     NEW.sale_id, 'INSERT', NULL,
                     row_to_json(NEW)::jsonb,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO sales_history (
            sale_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     NEW.sale_id, 'UPDATE',
                     row_to_json(OLD)::jsonb,
                     row_to_json(NEW)::jsonb,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO sales_history (
            sale_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     OLD.sale_id, 'DELETE',
                     row_to_json(OLD)::jsonb,
                     NULL,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Триггер для sales
CREATE TRIGGER trg_sales_history
    AFTER INSERT OR UPDATE OR DELETE ON sales
    FOR EACH ROW EXECUTE FUNCTION log_sales_changes();

-- Функция для логирования изменений в таблице vehicles
CREATE OR REPLACE FUNCTION log_vehicles_changes()
    RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO vehicles_history (
            vehicle_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     NEW.vehicle_id, 'INSERT', NULL,
                     row_to_json(NEW)::jsonb,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO vehicles_history (
            vehicle_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     NEW.vehicle_id, 'UPDATE',
                     row_to_json(OLD)::jsonb,
                     row_to_json(NEW)::jsonb,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO vehicles_history (
            vehicle_id, operation_type, old_value, new_value,
            username, hostname, application_name
        ) VALUES (
                     OLD.vehicle_id, 'DELETE',
                     row_to_json(OLD)::jsonb,
                     NULL,
                     CURRENT_USER, inet_client_addr()::VARCHAR,
                     current_setting('application_name', true)
                 );
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

-- Триггер для vehicles
CREATE TRIGGER trg_vehicles_history
    AFTER INSERT OR UPDATE OR DELETE ON vehicles
    FOR EACH ROW EXECUTE FUNCTION log_vehicles_changes();

-- Триггер для автоматической генерации номера контракта
CREATE OR REPLACE FUNCTION generate_contract_number()
    RETURNS TRIGGER AS $$
DECLARE
    v_year VARCHAR(4);
    v_month VARCHAR(2);
    v_sequence VARCHAR(6);
BEGIN
    IF NEW.contract_number IS NULL THEN
        v_year := TO_CHAR(CURRENT_DATE, 'YYYY');
        v_month := TO_CHAR(CURRENT_DATE, 'MM');
        v_sequence := LPAD(NEW.sale_id::TEXT, 6, '0');
        NEW.contract_number := 'AMK-' || v_year || v_month || '-' || v_sequence;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_generate_contract_number
    BEFORE INSERT ON sales
    FOR EACH ROW EXECUTE FUNCTION generate_contract_number();

-- Триггер для проверки резервирования техники при создании тест-драйва
CREATE OR REPLACE FUNCTION check_vehicle_for_test_drive()
    RETURNS TRIGGER AS $$
DECLARE
    v_vehicle_status VARCHAR(50);
    v_overlapping_count INTEGER;
BEGIN
    -- Проверка статуса техники
    SELECT status INTO v_vehicle_status
    FROM vehicles
    WHERE vehicle_id = NEW.vehicle_id;

    IF v_vehicle_status NOT IN ('В наличии', 'Зарезервировано') THEN
        RAISE EXCEPTION 'Техника недоступна для тест-драйва. Текущий статус: %', v_vehicle_status;
    END IF;

    -- Проверка на пересечение времени
    SELECT COUNT(*) INTO v_overlapping_count
    FROM test_drives
    WHERE vehicle_id = NEW.vehicle_id
      AND status = 'Запланирован'
      AND test_drive_id != COALESCE(NEW.test_drive_id, 0)
      AND (
        (NEW.scheduled_date, NEW.scheduled_date + (NEW.duration || ' minutes')::INTERVAL)
            OVERLAPS
        (scheduled_date, scheduled_date + (duration || ' minutes')::INTERVAL)
        );

    IF v_overlapping_count > 0 THEN
        RAISE EXCEPTION 'На это время уже запланирован другой тест-драйв';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_vehicle_test_drive
    BEFORE INSERT OR UPDATE ON test_drives
    FOR EACH ROW EXECUTE FUNCTION check_vehicle_for_test_drive();

-- Триггер для проверки остатков запчастей при добавлении в сервисный заказ
CREATE OR REPLACE FUNCTION check_spare_parts_stock()
    RETURNS TRIGGER AS $$
DECLARE
    v_available_quantity INTEGER;
    v_part_name VARCHAR(200);
BEGIN
    SELECT quantity_in_stock, part_name
    INTO v_available_quantity, v_part_name
    FROM spare_parts
    WHERE spare_part_id = NEW.spare_part_id;

    IF v_available_quantity < NEW.quantity THEN
        RAISE EXCEPTION 'Недостаточно запчастей "%". В наличии: %, требуется: %',
            v_part_name, v_available_quantity, NEW.quantity;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_check_spare_parts_stock
    BEFORE INSERT ON service_order_parts
    FOR EACH ROW EXECUTE FUNCTION check_spare_parts_stock();

-- Триггер для обновления остатков запчастей
CREATE OR REPLACE FUNCTION update_spare_parts_stock()
    RETURNS TRIGGER AS $$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        UPDATE spare_parts
        SET quantity_in_stock = quantity_in_stock - NEW.quantity
        WHERE spare_part_id = NEW.spare_part_id;
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        UPDATE spare_parts
        SET quantity_in_stock = quantity_in_stock + OLD.quantity - NEW.quantity
        WHERE spare_part_id = NEW.spare_part_id;
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        UPDATE spare_parts
        SET quantity_in_stock = quantity_in_stock + OLD.quantity
        WHERE spare_part_id = OLD.spare_part_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_spare_parts_stock
    AFTER INSERT OR UPDATE OR DELETE ON service_order_parts
    FOR EACH ROW EXECUTE FUNCTION update_spare_parts_stock();

-- Триггер для уведомления о низком остатке запчастей
CREATE OR REPLACE FUNCTION notify_low_spare_parts_stock()
    RETURNS TRIGGER AS $$
BEGIN
    IF NEW.quantity_in_stock <= NEW.min_quantity THEN
        -- Здесь можно добавить отправку уведомления
        RAISE NOTICE 'Низкий остаток запчастей: % (Артикул: %). Остаток: %, Минимум: %',
            NEW.part_name, NEW.part_number, NEW.quantity_in_stock, NEW.min_quantity;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_notify_low_stock
    AFTER UPDATE ON spare_parts
    FOR EACH ROW
    WHEN (NEW.quantity_in_stock <= NEW.min_quantity AND OLD.quantity_in_stock > OLD.min_quantity)
EXECUTE FUNCTION notify_low_spare_parts_stock();

-- Триггер для валидации продажи
CREATE OR REPLACE FUNCTION validate_sale()
    RETURNS TRIGGER AS $$
DECLARE
    v_vehicle_status VARCHAR(50);
BEGIN
    -- Проверка статуса техники
    SELECT status INTO v_vehicle_status
    FROM vehicles
    WHERE vehicle_id = NEW.vehicle_id;

    IF v_vehicle_status NOT IN ('В наличии', 'Зарезервировано') THEN
        RAISE EXCEPTION 'Невозможно продать технику со статусом: %', v_vehicle_status;
    END IF;

    -- Проверка что указан хотя бы один клиент
    IF NEW.customer_id IS NULL AND NEW.corporate_client_id IS NULL THEN
        RAISE EXCEPTION 'Необходимо указать клиента (физическое или юридическое лицо)';
    END IF;

    -- Проверка что не указаны оба типа клиентов
    IF NEW.customer_id IS NOT NULL AND NEW.corporate_client_id IS NOT NULL THEN
        RAISE EXCEPTION 'Нельзя указать одновременно физическое и юридическое лицо';
    END IF;

    -- Проверка корректности цен
    IF NEW.final_price > NEW.base_price THEN
        RAISE EXCEPTION 'Финальная цена не может быть больше базовой';
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_validate_sale
    BEFORE INSERT OR UPDATE ON sales
    FOR EACH ROW EXECUTE FUNCTION validate_sale();