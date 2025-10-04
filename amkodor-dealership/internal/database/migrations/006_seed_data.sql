-- Заполнение базы данных тестовыми данными

-- 1. Категории техники
INSERT INTO vehicle_categories (category_name, description) VALUES
                                                                ('Дорожно-строительная техника', 'Машины для строительства и ремонта дорог'),
                                                                ('Коммунальная техника', 'Техника для коммунального хозяйства'),
                                                                ('Погрузочная техника', 'Погрузчики и складская техника'),
                                                                ('Лесная техника', 'Машины для лесозаготовки'),
                                                                ('Сельскохозяйственная техника', 'Трактора и сельхозмашины');

-- 2. Типы техники
INSERT INTO vehicle_types (category_id, type_name, description) VALUES
                                                                    (1, 'Экскаваторы', 'Гусеничные и колесные экскаваторы'),
                                                                    (1, 'Бульдозеры', 'Гусеничные бульдозеры'),
                                                                    (1, 'Грейдеры', 'Автогрейдеры'),
                                                                    (2, 'Коммунальные машины', 'Уборочная техника'),
                                                                    (2, 'Мусоровозы', 'Машины для вывоза мусора'),
                                                                    (3, 'Фронтальные погрузчики', 'Колесные погрузчики'),
                                                                    (3, 'Вилочные погрузчики', 'Складские погрузчики'),
                                                                    (4, 'Трелевочные тракторы', 'Лесозаготовительная техника'),
                                                                    (5, 'Тракторы', 'Сельскохозяйственные тракторы');

-- 3. Производители
INSERT INTO manufacturers (manufacturer_name, country, website, contact_phone) VALUES
                                                                                   ('Амкодор', 'Беларусь', 'https://amkodor.by', '+375 17 200-00-00'),
                                                                                   ('МТЗ', 'Беларусь', 'https://mtz.by', '+375 17 370-00-00'),
                                                                                   ('БелАЗ', 'Беларусь', 'https://belaz.by', '+375 2175 300-00');

-- 4. Модели техники
INSERT INTO vehicle_models (model_name, type_id, manufacturer_id, description, specifications) VALUES
                                                                                                   ('Амкодор 332С', 1, 1, 'Универсальный экскаватор-погрузчик', '{"engine": "110 л.с.", "weight": "7500 кг", "bucket_capacity": "1.0 м³"}'),
                                                                                                   ('Амкодор 342С', 1, 1, 'Экскаватор-погрузчик повышенной мощности', '{"engine": "125 л.с.", "weight": "8200 кг", "bucket_capacity": "1.2 м³"}'),
                                                                                                   ('Амкодор ТО-18Б3', 2, 1, 'Трактор общего назначения', '{"engine": "180 л.с.", "weight": "9500 кг"}'),
                                                                                                   ('Амкодор 333ВL', 6, 1, 'Фронтальный погрузчик', '{"engine": "170 л.с.", "weight": "12000 кг", "bucket_capacity": "3.5 м³"}'),
                                                                                                   ('Амкодор 352С', 6, 1, 'Фронтальный погрузчик большой мощности', '{"engine": "220 л.с.", "weight": "15000 кг", "bucket_capacity": "5.0 м³"}'),
                                                                                                   ('МТЗ-82.1', 9, 2, 'Универсальный сельскохозяйственный трактор', '{"engine": "81 л.с.", "weight": "3545 кг"}'),
                                                                                                   ('МТЗ-920', 9, 2, 'Колесный трактор', '{"engine": "90 л.с.", "weight": "3850 кг"}'),
                                                                                                   ('Амкодор 12350', 4, 1, 'Коммунальная машина', '{"engine": "140 л.с.", "capacity": "8000 л"}'),
                                                                                                   ('Амкодор 6623', 7, 1, 'Вилочный погрузчик', '{"engine": "90 л.с.", "lift_capacity": "3000 кг"}'),
                                                                                                   ('Амкодор ТЛ-4', 8, 1, 'Трелевочный трактор', '{"engine": "210 л.с.", "weight": "16500 кг"}');

-- 5. Склады
INSERT INTO warehouses (warehouse_name, address, city, region, phone, manager_name, capacity) VALUES
                                                                                                  ('Главный склад Минск', 'ул. Промышленная, 1', 'Минск', 'Минская область', '+375 17 200-00-01', 'Иванов Иван Иванович', 100),
                                                                                                  ('Склад Гродно', 'ул. Складская, 15', 'Гродно', 'Гродненская область', '+375 152 123-45-67', 'Петров Петр Петрович', 50),
                                                                                                  ('Склад Витебск', 'пр. Московский, 88', 'Витебск', 'Витебская область', '+375 212 345-67-89', 'Сидоров Сидор Сидорович', 60),
                                                                                                  ('Склад Брест', 'ул. Московская, 250', 'Брест', 'Брестская область', '+375 162 234-56-78', 'Козлов Андрей Викторович', 40),
                                                                                                  ('Склад Гомель', 'ул. Речицкое шоссе, 100', 'Гомель', 'Гомельская область', '+375 232 456-78-90', 'Морозов Дмитрий Александрович', 55);

-- 6. Должности
INSERT INTO positions (position_name, base_salary, description) VALUES
                                                                    ('Менеджер по продажам', 1500.00, 'Продажа техники'),
                                                                    ('Старший менеджер по продажам', 2000.00, 'Продажа техники, управление отделом'),
                                                                    ('Механик', 1200.00, 'Ремонт и обслуживание техники'),
                                                                    ('Мастер сервисного центра', 1800.00, 'Руководство сервисным обслуживанием'),
                                                                    ('Кладовщик', 1000.00, 'Учет и хранение техники'),
                                                                    ('Директор филиала', 3000.00, 'Управление филиалом'),
                                                                    ('Бухгалтер', 1400.00, 'Ведение бухгалтерского учета'),
                                                                    ('Администратор', 2500.00, 'Общее управление');

-- 7. Сотрудники (пароль для всех: password123, хеш bcrypt)
INSERT INTO employees (first_name, last_name, middle_name, position_id, warehouse_id, email, phone, password_hash, salary) VALUES
                                                                                                                               ('Иван', 'Иванов', 'Иванович', 8, 1, 'admin@amkodor.by', '+375 29 111-11-11', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 3000.00),
                                                                                                                               ('Алексей', 'Петров', 'Сергеевич', 1, 1, 'petrov@amkodor.by', '+375 29 222-22-22', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 1800.00),
                                                                                                                               ('Мария', 'Сидорова', 'Андреевна', 1, 1, 'sidorova@amkodor.by', '+375 29 333-33-33', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 1700.00),
                                                                                                                               ('Дмитрий', 'Козлов', 'Викторович', 3, 1, 'kozlov@amkodor.by', '+375 29 444-44-44', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 1200.00),
                                                                                                                               ('Ольга', 'Морозова', 'Петровна', 2, 2, 'morozova@amkodor.by', '+375 29 555-55-55', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 2200.00),
                                                                                                                               ('Андрей', 'Волков', 'Николаевич', 1, 3, 'volkov@amkodor.by', '+375 29 666-66-66', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 1600.00),
                                                                                                                               ('Елена', 'Новикова', 'Владимировна', 4, 1, 'novikova@amkodor.by', '+375 29 777-77-77', '$2a$10$X.J9H6rMVUZ5YdLrN3vEh.OHLxQE6.YxQXXvhW0l8bWAh7JV0YnOi', 1900.00);

-- 8. Техника в наличии
INSERT INTO vehicles (model_id, warehouse_id, vin, serial_number, manufacture_year, color, price, discount, status) VALUES
                                                                                                                        (1, 1, 'AMK332S2023001', 'AMK001', 2023, 'Желтый', 85000.00, 0, 'В наличии'),
                                                                                                                        (1, 1, 'AMK332S2023002', 'AMK002', 2023, 'Желтый', 85000.00, 5, 'В наличии'),
                                                                                                                        (2, 1, 'AMK342S2023001', 'AMK003', 2023, 'Оранжевый', 95000.00, 0, 'В наличии'),
                                                                                                                        (3, 1, 'AMKTO18B2023001', 'AMK004', 2023, 'Красный', 120000.00, 0, 'В наличии'),
                                                                                                                        (4, 2, 'AMK333VL2023001', 'AMK005', 2023, 'Желтый', 150000.00, 0, 'В наличии'),
                                                                                                                        (4, 2, 'AMK333VL2023002', 'AMK006', 2023, 'Желтый', 150000.00, 3, 'В наличии'),
                                                                                                                        (5, 1, 'AMK352S2023001', 'AMK007', 2023, 'Желтый', 180000.00, 0, 'В наличии'),
                                                                                                                        (6, 3, 'MTZ821-2023001', 'MTZ001', 2023, 'Красный', 45000.00, 0, 'В наличии'),
                                                                                                                        (6, 3, 'MTZ821-2023002', 'MTZ002', 2023, 'Красный', 45000.00, 5, 'В наличии'),
                                                                                                                        (7, 3, 'MTZ920-2023001', 'MTZ003', 2023, 'Красный', 52000.00, 0, 'В наличии'),
                                                                                                                        (8, 1, 'AMK12350-2023001', 'AMK008', 2023, 'Белый', 110000.00, 0, 'В наличии'),
                                                                                                                        (9, 1, 'AMK6623-2023001', 'AMK009', 2023, 'Оранжевый', 65000.00, 0, 'В наличии'),
                                                                                                                        (10, 2, 'AMKTL4-2023001', 'AMK010', 2023, 'Зеленый', 165000.00, 0, 'В наличии'),
                                                                                                                        (1, 4, 'AMK332S2023003', 'AMK011', 2023, 'Желтый', 85000.00, 0, 'Зарезервировано'),
                                                                                                                        (2, 5, 'AMK342S2023002', 'AMK012', 2023, 'Оранжевый', 95000.00, 0, 'В наличии');

-- 9. Клиенты
INSERT INTO customers (first_name, last_name, middle_name, phone, email, address, discount_percent, is_vip) VALUES
                                                                                                                ('Александр', 'Смирнов', 'Петрович', '+375 29 123-45-67', 'smirnov@mail.ru', 'г. Минск, ул. Ленина, 10', 0, false),
                                                                                                                ('Екатерина', 'Кузнецова', 'Ивановна', '+375 29 234-56-78', 'kuznetsova@gmail.com', 'г. Гродно, ул. Советская, 25', 5, false),
                                                                                                                ('Николай', 'Попов', 'Александрович', '+375 29 345-67-89', 'popov@tut.by', 'г. Витебск, пр. Победителей, 50', 10, true),
                                                                                                                ('Светлана', 'Федорова', 'Сергеевна', '+375 29 456-78-90', 'fedorova@yandex.ru', 'г. Брест, ул. Московская, 15', 0, false),
                                                                                                                ('Владимир', 'Соколов', 'Михайлович', '+375 29 567-89-01', 'sokolov@mail.ru', 'г. Гомель, ул. Кирова, 30', 15, true);

-- 10. Корпоративные клиенты
INSERT INTO corporate_clients (company_name, tax_id, legal_address, contact_person, phone, email, discount_percent, contract_number) VALUES
                                                                                                                                         ('ООО "СтройТех"', '123456789', 'г. Минск, ул. Промышленная, 100', 'Иванов И.И.', '+375 17 200-10-10', 'info@stroytech.by', 10, 'CT-2023-001'),
                                                                                                                                         ('ЗАО "АгроМаш"', '987654321', 'г. Гродно, ул. Фабричная, 50', 'Петров П.П.', '+375 152 300-20-20', 'contact@agromash.by', 15, 'CT-2023-002'),
                                                                                                                                         ('ОАО "Лесхоз"', '456789123', 'г. Витебск, ул. Лесная, 25', 'Сидоров С.С.', '+375 212 400-30-30', 'office@leshoz.by', 12, 'CT-2023-003'),
                                                                                                                                         ('ИП Козлов А.В.', '789123456', 'г. Брест, ул. Коммунальная, 10', 'Козлов А.В.', '+375 162 500-40-40', 'kozlov.av@mail.ru', 5, 'CT-2023-004');

-- 11. Продажи
INSERT INTO sales (vehicle_id, customer_id, employee_id, sale_date, base_price, discount_amount, final_price, payment_type, status) VALUES
                                                                                                                                        (14, 3, 2, CURRENT_DATE - INTERVAL '5 days', 85000.00, 8500.00, 76500.00, 'Безналичный', 'Завершена'),
                                                                                                                                        (NULL, NULL, 2, CURRENT_DATE - INTERVAL '3 days', 150000.00, 15000.00, 135000.00, 'Кредит', 'Завершена');

-- Обновляем статус проданной техники
UPDATE vehicles SET status = 'Продано' WHERE vehicle_id = 14;

-- 12. Тест-драйвы
INSERT INTO test_drives (vehicle_id, customer_id, employee_id, scheduled_date, duration, status) VALUES
                                                                                                     (1, 1, 2, CURRENT_TIMESTAMP + INTERVAL '2 days', 60, 'Запланирован'),
                                                                                                     (3, 2, 3, CURRENT_TIMESTAMP + INTERVAL '3 days', 90, 'Запланирован'),
                                                                                                     (7, NULL, 2, CURRENT_TIMESTAMP + INTERVAL '5 days', 120, 'Запланирован');

-- 13. Запасные части
INSERT INTO spare_parts (part_number, part_name, model_id, price, quantity_in_stock, min_quantity, warehouse_id) VALUES
                                                                                                                     ('AMK-FLT-001', 'Фильтр масляный', 1, 45.00, 50, 10, 1),
                                                                                                                     ('AMK-FLT-002', 'Фильтр воздушный', 1, 65.00, 30, 8, 1),
                                                                                                                     ('AMK-BRK-001', 'Тормозные колодки', 1, 120.00, 25, 5, 1),
                                                                                                                     ('AMK-HYD-001', 'Гидравлический шланг', NULL, 85.00, 40, 10, 1),
                                                                                                                     ('MTZ-FLT-001', 'Фильтр топливный МТЗ', 6, 38.00, 60, 15, 3),
                                                                                                                     ('AMK-TYR-001', 'Шина 18.4-30', NULL, 850.00, 20, 5, 1),
                                                                                                                     ('AMK-BTR-001', 'Аккумулятор 12V 190Ah', NULL, 320.00, 15, 3, 1);

-- 14. Сервисные заказы
INSERT INTO service_orders (vehicle_id, customer_id, employee_id, order_date, service_type, description, cost, status) VALUES
                                                                                                                           (1, 1, 4, CURRENT_DATE - INTERVAL '10 days', 'ТО-1', 'Первое техническое обслуживание', 450.00, 'Завершен'),
                                                                                                                           (8, 3, 4, CURRENT_DATE - INTERVAL '5 days', 'Ремонт', 'Замена масла и фильтров', 280.00, 'Завершен'),
                                                                                                                           (3, 2, 4, CURRENT_DATE, 'ТО-2', 'Второе техническое обслуживание', 0, 'В работе');

-- 15. Запчасти для сервисных заказов
INSERT INTO service_order_parts (service_order_id, spare_part_id, quantity, unit_price) VALUES
                                                                                            (1, 1, 2, 45.00),
                                                                                            (1, 2, 1, 65.00),
                                                                                            (2, 1, 1, 45.00),
                                                                                            (2, 5, 1, 38.00);

-- 16. Поставки
INSERT INTO supplies (manufacturer_id, warehouse_id, supply_date, expected_arrival_date, status, invoice_number) VALUES
                                                                                                                     (1, 1, CURRENT_DATE - INTERVAL '30 days', CURRENT_DATE - INTERVAL '25 days', 'Получено', 'INV-2023-001'),
                                                                                                                     (2, 3, CURRENT_DATE - INTERVAL '20 days', CURRENT_DATE - INTERVAL '15 days', 'Получено', 'INV-2023-002'),
                                                                                                                     (1, 2, CURRENT_DATE, CURRENT_DATE + INTERVAL '10 days', 'В пути', 'INV-2023-003');

-- 17. Элементы поставок
INSERT INTO supply_items (supply_id, model_id, quantity, unit_price) VALUES
                                                                         (1, 1, 5, 82000.00),
                                                                         (1, 2, 3, 92000.00),
                                                                         (2, 6, 10, 43000.00),
                                                                         (3, 4, 4, 145000.00);

-- 18. Конфигурация SSIS
INSERT INTO ssis_configuration (config_name, config_value, description) VALUES
                                                                            ('report_start_date', '2023-01-01', 'Начальная дата для отчетов'),
                                                                            ('report_end_date', '2023-12-31', 'Конечная дата для отчетов'),
                                                                            ('export_path', '/exports/', 'Путь для экспорта отчетов'),
                                                                            ('email_notifications', 'true', 'Включить email уведомления');

-- Логирование выполнения
INSERT INTO ssis_log (package_name, execution_date, rows_processed, status, username, execution_duration) VALUES
                                                                                                              ('SalesExport', CURRENT_TIMESTAMP - INTERVAL '1 day', 150, 'Success', 'admin', 45),
                                                                                                              ('InventoryExport', CURRENT_TIMESTAMP - INTERVAL '2 days', 200, 'Success', 'admin', 60);

COMMIT;