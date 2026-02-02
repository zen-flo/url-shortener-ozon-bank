# url-shortener
Сервис для сокращения URL с двумя типами хранилища:
- in-memory
- PostgreSQL.
---

> Данный проект является решением тестового задания от Ozon Банка на позицию стажёра-разработчика Golang.

---
## Описание
Сервис предоставляет API для создания коротких ссылок и редиректа по ним.

### Ключевые особенности
- Ссылка уникальна: один оригинальный URL – одна короткая ссылка.
- Длина короткой ссылки – 10 символов.
- Короткая ссылка состоит из:
  - символов латинского алфавита в нижнем и верхнем регистре (A–Z, a–z)
  - цифр (0–9)
  - символа нижнего подчёркивания (_).
- HTTP запросы:
  - **POST `/shorten`** – сохраняет оригинальный URL и возвращает сокращённый
  - **GET `/{short_code}`** – редиректит на оригинальный URL.
- Сборка и запуск выполняются через Docker (подробнее в разделе _"Настройка и сборка"_)
- Две реализации хранилища:
  - в памяти приложения (`memory`)
  - PostgreSQL (`postgres`).
- Конфигурация через переменные окружения:
  - `STORE_TYPE` - выбор хранилища (по умолчанию: `memory`)
  - `PORT` – выбор порта (по умолчанию: `8080`)
  - `POSTGRES_DSN` – строка с параметрами подключения к БД (при выборе типа хранилища `postgres`)
  - `POSTGRES_TEST_DSN` – строка с параметрами подключения к тестовой БД (при запуске unit-тестов для хранилища `postgres`).
- Приложение покрыто unit-тестами (`memory`, `postgres`, `service`).

---
## Структура проекта
```shell
.
├── Dockerfile # сборка Docker-образа сервиса
├── README.md 
├── cmd
│   └── server # основной исполняемый файл
│       └── main.go
├── go.mod # модули Go
├── go.sum
└── internal
    ├── config # конфигурация (порт, тип хранилища)
    │   └── config.go
    ├── handler # HTTP-обработчики
    │   ├── handler.go
    │   ├── redirect.go
    │   └── shorten.go
    ├── service # бизнес-логика
    │   ├── shortener.go
    │   └── shortener_test.go
    ├── shortcode
    │   └── codegen.go
    └── storage
        ├── memory # реализация in-memory хранилища
        │   ├── memory.go
        │   └── memory_test.go
        ├── postgres # реализация хранилища PostgreSQL
        │   ├── postgres.go
        │   └── postgres_test.go
        └── storage.go
```

---
## Настройка и сборка
### 1. Клонирование репозитория
```shell
git clone https://github.com/zen-flo/url-shortener-ozon-bank.git
cd url-shortener-ozon-bank
```
### 2. Сборка Docker-образа
```shell
docker build -t url-shortener .
```
### 3. Запуск сервиса
#### 3.1. Запуск с PostgreSQL
##### 3.1.1. Запуск БД PostgreSQL
```shell
docker run -d \
  --name shortener-db \
  -e POSTGRES_USER=shortener_user \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=shortener \
  -p 5432:5432 \
  postgres:18.1
```
##### 3.1.2. Запуск сервиса
```shell
docker run --rm -p 8080:8080 \
  -e STORE_TYPE=postgres \
  -e POSTGRES_DSN="postgres://shortener_user:secret@host.docker.internal:5432/shortener?sslmode=disable" \
  url-shortener
```
#### 3.2. Запуск с in-memory
```shell
docker run --rm -p 8080:8080 \
  -e STORE_TYPE=memory \
  url-shortener
```

---
## Примеры HTTP-запросов
### Проверка состояния сервиса
```shell
curl http://localhost:8080/health
```
#### Ответ:
```shell
  OK
```
### Создание короткой ссылки
```shell
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url":"https://example.com"}'
```
#### Ответ:
```JSON
{
  "short_code":"23_nn35cw8"
}
```
### Редирект по короткой ссылке
```shell
curl -v http://localhost:8080/23_nn35cw8
```
#### Ответ:
Редиректит (`HTTP 302`) на `https://example.com`.

---
## Алгоритм генерации коротких ссылок
Короткая ссылка генерируется в виде строки длиной **10 символов** и состоит из следующего алфавита:
- латинские буквы в нижнем и верхнем регистре (A–Z, a–z)
- цифры (0–9)
- символ нижнего подчёркивания (_).

Для генерации ссылок используется криптоустойчивый генератор случайных данных из пакета `crypto/rand`.

При получении запроса на сокращение URL сначала проверяется,
   существует ли данный оригинальный URL в хранилище.
- Если URL уже сохранён, сервис возвращает ранее сгенерированный короткий код. 
- Если URL новый, генерируется случайная строка длиной 10 символов
   из допустимого алфавита. 

Для обеспечения уникальности:
- в `in-memory` хранилище проверяется отсутствие такого кода в мапе
- в PostgreSQL используется уникальный индекс и повторная генерация
  кода при коллизии. 

Сгенерированный короткий код сохраняется вместе с оригинальным URL
   и возвращается клиенту.

---
## Запуск тестов
### Тесты PostgreSQL
#### 1. Запуск тестовой БД (если не запущена)
```shell
docker run -d \
  --name shortener-test-db \
  -e POSTGRES_USER=shortener_user \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=shortener_test \
  -p 5433:5432 \
  postgres:18.1
```
#### 2. Запуск тестов
```shell
export POSTGRES_TEST_DSN="postgres://shortener_user:secret@localhost:5433/shortener_test?sslmode=disable"
go test ./internal/storage/postgres -v
```
### Тесты Go
```shell
# Для всех пакетов
go test -v ./... 

# С покрытием
go test ./... -cover
```

> `service` – 100.0%
> 
> `memory` – 90.0%
> 
> `postgres` – 63.9%

---