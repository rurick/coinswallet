## Структура проекта:

* /build/ - содержит инструменты для запуска тестов и окружения
* /cmd/ - приложение http-сервер обеспечивающий прием и обработку REST API запросов
* /doc/ - документация к проекту
* /internal/ - реализация микросервиса. подробнее в файле architecture.md
* /pkg/ - пакеты-утилиты для работы сервиса
* /.env - настройки соединения с БД и настройки кеширования

## Тестирование
build/test.sh Для запуска unit-тестов и интеграционных тестов   