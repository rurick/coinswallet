## Структура проекта:

* /build/ - содержит инструменты для запуска тестов и окружения
* /cmd/ - приложение http-сервер обеспечивающий прием и обработку REST API запросов
* /doc/ - документация к проекту
* /internal/ - реализация микросервиса. подробнее в файле architecture.md
* /pkg/ - пакеты-утилиты для работы сервиса
* /.env - настройки соединения с БД и настройки кеширования

### /build/
* /build/http/ - здесь находятся файлы с http запросами для ручного тестирования
* init.sql - скрипт инициализации БД
* docker-compose.yml - запуск микросервиса (вместе с субд)
* pgdocker_up.sh - запуск субд postgres в докере
* pgdocker_init.sh - инициализация БД
* pgdocker_down.sh - остановка субд postgres в докере
* test.sh - запуск unit-тестов и интеграционных тестов 
* test_api.sh - запуск автоматического тестирования api

## Тестирование
```shell
$ cd build
$ sudo ./pgdocker_up.sh
$ sudo ./pgdocker_init.sh
$ ./test.sh
$ sudo ./pgdocker_down.sh

$ sudo docker-compose up -d
$ ./test_api.sh
$ sudo docker-compose down
```
## Конфигурация 
Конфигурирование осуществляется через переменные окружения.

Настраивается время жизни кэша и параметры соединения с СУБД:
```shell
# PostgreSQL connection
PGSQL_HOST=127.0.0.1
PGSQL_NAME=coins
PGSQL_USER=coins
PGSQL_PASS=coins
PGSQL_PORT=5433
# Memory cache settings (in minutes)
CacheExpTime=10
```

## Запуск приложения
Для запуска приложения с использованием docker-compose: build/docker-compose.yml
```shell
cd build/
sudo docker-compose up -d
```

`После запуска приложения можно воспользоваться ручным тестирование выполняя запросы из build/http/test.http`
