-- Представления для автосалона Амкодор (PostgreSQL)

-- 1. Полная информация о технике
CREATE OR REPLACE VIEW vw_vehicles_full_info AS
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
    ROUND(v.price * (1 - v.discount / 100), 2) AS final_price,
    v.status,
    w.warehouse_name,
    w.city AS warehouse_city,
    v.arrival_date,
    v.created_at,
    vm.description,
    vm.specifications
FROM vehicles v
         INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
         INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
         INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
         INNER JOIN manufacturers m ON vm.manufacturer_id = m.manufacturer_id
         INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id;

-- 2. Полная информация о продажах
CREATE OR REPLACE VIEW vw_sales_full_info AS
SELECT
    s.sale_id,
    s.contract_number,
    s.sale_date,
    v.vin,
    vm.model_name,
    vt.type_name,
    CASE
        WHEN s.customer_id IS NOT NULL THEN c.last_name || ' ' || c.first_name
        ELSE cc.company_name
        END AS client_name,
    CASE
        WHEN s.customer_id IS NOT NULL THEN c.phone
        ELSE cc.phone
        END AS client_phone,
    CASE
        WHEN s.customer_id IS NOT NULL THEN 'Физ. лицо'
        ELSE 'Юр. лицо'
        END AS client_type,
    e.last_name || ' ' || e.first_name AS manager_name,
    p.position_name,
    s.base_price,
    s.discount_amount,
    s.final_price,
    s.payment_type,
    s.status,
    w.warehouse_name,
    w.city AS warehouse_city
FROM sales s
         INNER JOIN vehicles v ON s.vehicle_id = v.vehicle_id
         INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
         INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
         INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
         INNER JOIN employees e ON s.employee_id = e.employee_id
         INNER JOIN positions p ON e.position_id = p.position_id
         LEFT JOIN customers c ON s.customer_id = c.customer_id
         LEFT JOIN corporate_clients cc ON s.corporate_client_id = cc.corporate_client_id;

-- 3. Доступная техника для продажи
CREATE OR REPLACE VIEW vw_available_vehicles AS
SELECT
    v.vehicle_id,
    vm.model_name,
    vt.type_name,
    vc.category_name,
    m.manufacturer_name,
    v.manufacture_year,
    v.color,
    v.price,
    v.discount,
    ROUND(v.price * (1 - v.discount / 100), 2) AS final_price,
    w.warehouse_name,
    w.city,
    w.phone AS warehouse_phone,
    vm.description,
    vm.specifications,
    CURRENT_DATE - v.arrival_date AS days_in_stock
FROM vehicles v
         INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
         INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
         INNER JOIN vehicle_categories vc ON vt.category_id = vc.category_id
         INNER JOIN manufacturers m ON vm.manufacturer_id = m.manufacturer_id
         INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
WHERE v.status = 'В наличии' AND w.is_active = TRUE;

-- 4. Полная информация о сотрудниках
CREATE OR REPLACE VIEW vw_employees_full_info AS
SELECT
    e.employee_id,
    e.last_name || ' ' || e.first_name || COALESCE(' ' || e.middle_name, '') AS full_name,
    e.first_name,
    e.last_name,
    e.middle_name,
    p.position_name,
    p.base_salary,
    e.salary,
    w.warehouse_name,
    w.city AS warehouse_city,
    e.email,
    e.phone,
    e.hire_date,
    EXTRACT(YEAR FROM AGE(CURRENT_DATE, e.hire_date)) AS years_of_service,
    e.is_active
FROM employees e
         INNER JOIN positions p ON e.position_id = p.position_id
         LEFT JOIN warehouses w ON e.warehouse_id = w.warehouse_id;

-- 5. Полная информация о тест-драйвах
CREATE OR REPLACE VIEW vw_test_drives_full_info AS
SELECT
    td.test_drive_id,
    td.scheduled_date,
    td.duration,
    td.status,
    vm.model_name,
    vt.type_name,
    v.color,
    v.manufacture_year,
    CASE
        WHEN td.customer_id IS NOT NULL THEN c.last_name || ' ' || c.first_name
        ELSE cc.company_name
        END AS client_name,
    CASE
        WHEN td.customer_id IS NOT NULL THEN c.phone
        ELSE cc.phone
        END AS client_phone,
    e.last_name || ' ' || e.first_name AS manager_name,
    e.phone AS manager_phone,
    td.feedback_rating,
    td.feedback_comment,
    w.warehouse_name,
    w.city AS warehouse_city
FROM test_drives td
         INNER JOIN vehicles v ON td.vehicle_id = v.vehicle_id
         INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
         INNER JOIN vehicle_types vt ON vm.type_id = vt.type_id
         INNER JOIN warehouses w ON v.warehouse_id = w.warehouse_id
         INNER JOIN employees e ON td.employee_id = e.employee_id
         LEFT JOIN customers c ON td.customer_id = c.customer_id
         LEFT JOIN corporate_clients cc ON td.corporate_client_id = cc.corporate_client_id;

-- 6. Полная информация о сервисных заказах
CREATE OR REPLACE VIEW vw_service_orders_full_info AS
SELECT
    so.service_order_id,
    so.order_date,
    so.completion_date,
    so.service_type,
    so.status,
    so.cost,
    vm.model_name,
    v.vin,
    CASE
        WHEN so.customer_id IS NOT NULL THEN c.last_name || ' ' || c.first_name
        ELSE cc.company_name
        END AS client_name,
    CASE
        WHEN so.customer_id IS NOT NULL THEN c.phone
        ELSE cc.phone
        END AS client_phone,
    e.last_name || ' ' || e.first_name AS master_name,
    (SELECT COALESCE(SUM(sop.quantity * sop.unit_price), 0)
     FROM service_order_parts sop
     WHERE sop.service_order_id = so.service_order_id) AS parts_cost,
    so.description
FROM service_orders so
         INNER JOIN vehicles v ON so.vehicle_id = v.vehicle_id
         INNER JOIN vehicle_models vm ON v.model_id = vm.model_id
         INNER JOIN employees e ON so.employee_id = e.employee_id
         LEFT JOIN customers c ON so.customer_id = c.customer_id
         LEFT JOIN corporate_clients cc ON so.corporate_client_id = cc.corporate_client_id;

-- 7. Запасные части с остатками
CREATE OR REPLACE VIEW vw_spare_parts_inventory AS
SELECT
    sp.spare_part_id,
    sp.part_number,
    sp.part_name,
    vm.model_name,
    sp.price,
    sp.quantity_in_stock,
    sp.min_quantity,
    CASE
        WHEN sp.quantity_in_stock = 0 THEN 'Нет в наличии'
        WHEN sp.quantity_in_stock <= sp.min_quantity THEN 'Требуется заказ'
        WHEN sp.quantity_in_stock <= sp.min_quantity * 2 THEN 'Низкий остаток'
        ELSE 'В наличии'
        END AS stock_status,
    w.warehouse_name,
    w.city AS warehouse_city
FROM spare_parts sp
         LEFT JOIN vehicle_models vm ON sp.model_id = vm.model_id
         INNER JOIN warehouses w ON sp.warehouse_id = w.warehouse_id;

-- 8. Статистика продаж по менеджерам
CREATE OR REPLACE VIEW vw_sales_statistics_by_manager AS
SELECT
    e.employee_id,
    e.last_name || ' ' || e.first_name AS manager_name,
    w.warehouse_name,
    COUNT(s.sale_id) AS total_sales,
    COALESCE(SUM(s.final_price), 0) AS total_revenue,
    COALESCE(AVG(s.final_price), 0) AS average_sale_price,
    MIN(s.sale_date) AS first_sale_date,
    MAX(s.sale_date) AS last_sale_date
FROM employees e
         LEFT JOIN sales s ON e.employee_id = s.employee_id AND s.status = 'Завершена'
         LEFT JOIN warehouses w ON e.warehouse_id = w.warehouse_id
         INNER JOIN positions p ON e.position_id = p.position_id
WHERE p.position_name ILIKE '%менеджер%' OR p.position_name ILIKE '%продаж%'
GROUP BY e.employee_id, e.last_name, e.first_name, w.warehouse_name;

-- 9. Статистика продаж по категориям
CREATE OR REPLACE VIEW vw_sales_statistics_by_category AS
SELECT
    vc.category_id,
    vc.category_name,
    COUNT(s.sale_id) AS total_sales,
    COALESCE(SUM(s.final_price), 0) AS total_revenue,
    COALESCE(AVG(s.final_price), 0) AS average_sale_price,
    MIN(s.sale_date) AS first_sale_date,
    MAX(s.sale_date) AS last_sale_date
FROM vehicle_categories vc
         INNER JOIN vehicle_types vt ON vc.category_id = vt.category_id
         INNER JOIN vehicle_models vm ON vt.type_id = vm.type_id
         INNER JOIN vehicles v ON vm.model_id = v.model_id
         LEFT JOIN sales s ON v.vehicle_id = s.vehicle_id AND s.status = 'Завершена'
GROUP BY vc.category_id, vc.category_name;

-- 10. Полная информация о поставках
CREATE OR REPLACE VIEW vw_supplies_full_info AS
SELECT
    su.supply_id,
    su.supply_date,
    su.expected_arrival_date,
    su.actual_arrival_date,
    su.status,
    su.total_cost,
    su.invoice_number,
    m.manufacturer_name,
    m.country,
    w.warehouse_name,
    w.city AS warehouse_city,
    (SELECT COUNT(*) FROM supply_items si WHERE si.supply_id = su.supply_id) AS items_count,
    (SELECT COALESCE(SUM(si.quantity), 0) FROM supply_items si WHERE si.supply_id = su.supply_id) AS total_quantity
FROM supplies su
         INNER JOIN manufacturers m ON su.manufacturer_id = m.manufacturer_id
         INNER JOIN warehouses w ON su.warehouse_id = w.warehouse_id;

-- 11. Все клиенты (физ. и юр. лица)
CREATE OR REPLACE VIEW vw_all_clients AS
SELECT
    'CUSTOMER' AS client_type,
    customer_id AS client_id,
    last_name || ' ' || first_name AS client_name,
    phone,
    email,
    discount_percent,
    CASE WHEN is_vip THEN 'VIP' ELSE 'Обычный' END AS client_category,
    created_at
FROM customers
UNION ALL
SELECT
    'CORPORATE' AS client_type,
    corporate_client_id AS client_id,
    company_name AS client_name,
    phone,
    email,
    discount_percent,
    'Корпоративный' AS client_category,
    created_at
FROM corporate_clients;

-- 12. Статистика по складам
CREATE OR REPLACE VIEW vw_warehouses_statistics AS
SELECT
    w.warehouse_id,
    w.warehouse_name,
    w.city,
    w.region,
    w.capacity,
    COUNT(DISTINCT v.vehicle_id) AS vehicles_in_stock,
    COUNT(DISTINCT CASE WHEN v.status = 'В наличии' THEN v.vehicle_id END) AS available_vehicles,
    COUNT(DISTINCT e.employee_id) AS employees_count,
    COALESCE(SUM(ROUND(v.price * (1 - v.discount / 100), 2)), 0) AS total_inventory_value
FROM warehouses w
         LEFT JOIN vehicles v ON w.warehouse_id = v.warehouse_id
         LEFT JOIN employees e ON w.warehouse_id = e.warehouse_id AND e.is_active = TRUE
WHERE w.is_active = TRUE
GROUP BY w.warehouse_id, w.warehouse_name, w.city, w.region, w.capacity;

-- 13. Дашборд - общая статистика
CREATE OR REPLACE VIEW vw_dashboard_statistics AS
SELECT
    (SELECT COUNT(*) FROM vehicles WHERE status = 'В наличии') AS available_vehicles,
    (SELECT COUNT(*) FROM sales WHERE sale_date >= CURRENT_DATE - INTERVAL '30 days' AND status = 'Завершена') AS sales_last_month,
    (SELECT COALESCE(SUM(final_price), 0) FROM sales WHERE sale_date >= CURRENT_DATE - INTERVAL '30 days' AND status = 'Завершена') AS revenue_last_month,
    (SELECT COUNT(*) FROM customers) AS total_customers,
    (SELECT COUNT(*) FROM corporate_clients) AS total_corporate_clients,
    (SELECT COUNT(*) FROM test_drives WHERE scheduled_date >= CURRENT_TIMESTAMP AND status = 'Запланирован') AS upcoming_test_drives,
    (SELECT COUNT(*) FROM service_orders WHERE status = 'В работе') AS active_service_orders;