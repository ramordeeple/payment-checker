# 🏦 Payment Checker Service

### Сервис валидации платежей с использованием мок-сервиса ЦБ.

* Проверяет платежи в рублях и зарубежных валютах.
* Валюта и курс берутся из локальной базы через мок `/cbr`.
* Интеграция через HTTP и gRPC для микросервисного взаимодействия.
* Swagger документация доступна для удобного тестирования.

---

## ⚙️ Стек технологий

| Компонент       | Технологии           |
|-----------------|----------------------|
| Backend         | Go, net/http, gRPC   |
| Database        | PostgreSQL           |
| API Docs        | Swagger              |
| Dev Environment | Docker Compose       |
| Testing         | Go tests, .http file |
| Build           | Go modules, Makefile |

---

## 🧱 Архитектура проекта

```
cmd/payment-checker — точка входа (main.go)
internal/
  adapter/          — HTTP/gRPC handlers
  domain/           — DDD ядро
  usecase/          — бизнес-логика, Policy, Converter, тесты
  repository/       — работа с базой 
  port/             — интерфейсы 
docs/               — сгенерированная Swagger документация
```

---

## 🚀 Быстрый старт

### 1️⃣ Предварительные требования

* Go 1.21+
* Docker Desktop  (для базы данных)

---

### 2️⃣ Поднять базу данных

Пример `docker-compose.yml`:

```yaml
services:
  db:
    image: postgres:latest
    command: ["postgres", "-c", "port=5433"]
    environment:
      POSTGRES_DB: payment_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 123
    ports:
      - "5433:5433"
    volumes:
      - ./migrations:/docker-entrypoint-initdb.d
```

В корне проекта:

```bash
docker compose up -d
```

Будут применены указанные ниже параметры значений:

| Параметр | Значение   |
|----------|------------|
| Host     | localhost  |
| Port     | 5433       |
| Database | payment_db |
| Username | postgres   |
| Password | 123        |

Либо же правой кнопкой мыши по `docker-compose.yml` и `Run 'payment_db'`, если в **GoLand**.

> ⚠️ Порт 5433 выбран, чтобы не конфликтовать с локальным 5432.


---

### 3️⃣ Настройка переменных окружения

Указан в `main.go` только для базы данных:

```
DB_URL=postgres://postgres:123@localhost:5433/payment_db?sslmode=disable
```

---

### 4️⃣ Запуск приложения

Через Makefile:

```bash
make run
```

или напрямую:

```bash
cd cmd/payment-checker
go run main.go
```

После запуска:

* HTTP API: `http://localhost:8080`
* gRPC сервер: `localhost:9090`
* Swagger UI: `http://localhost:8080/swagger/index.html`

---

## 🌐 API

### Проверка 

```http
POST http://localhost:8080/validate
Content-Type: application/json

{
  "provider": "provider",
  "amount": 751234,
  "currency": "RUB",
  "date": "2025-10-10"
}
```

 Получение курса валюты (мок ЦБ)

```http
GET http://localhost:8080/cbr?currency=USD&date=2025-10-21
Accept: application/xml
```

**Пример ответа:**

```xml
<Valute>
<CharCode>USD</CharCode>
<Nominal>1</Nominal>
<Value>751234</Value>
</Valute>
```

> ⚠️ Здесь `Value` имеет тип `int64` во избежания ошибок с плавающей точкой. По сути **USD** получен в центах.
---

## 📡 gRPC

### Зачем gRPC

* Позволяет микросервисам общаться напрямую с **proto buffers**, что важно для интеграционных e2e тестов.

* Обеспечивает типизированный контракт через **.proto**.

* Подходит для параллельных и асинхронных вызовов.

Пример **gRPC** сервиса:

```proto
service PaymentChecker {
  rpc Validate(ValidatePaymentRequest) returns (ValidatePaymentResponse);
}
```

---

## 🪢 Реализация мок-сервиса

* `/cbr` — обращается к базе `rates` и `currencies`.
* **Гибкость:** данные можно менять в базе без изменения кода.
* **Воспроизводимость:** одинаковые запросы всегда возвращают одинаковый результат.
* **Надёжность:** тесты не зависят от внешних API и сети.

### Микросервисная архитектура:

* Параметры и заголовки нельзя использовать для управления моками.
* Единый тестовый стенд запускает сотни e2e тестов одновременно.
* Детермированные данные из базы обеспечивают стабильность и идемпотентность.

Мок-сервис `/cbr` обращается к локальной базе (`currencies` и `rates`) и возвращает
курс валюты на заданную дату в `JSON` или `xml`.

### Это обеспечивает:

* Гибкость — можно менять данные в базе без изменения кода.

* Воспроизводимость — одинаковый запрос всегда даёт одинаковый результат.

* Надёжность — не зависит от внешнего API или сети.

* В микросервисной архитектуре с параллельными e2e тестами нельзя управлять моками через параметры или заголовки. 
Данные из базы позволяют тестам быть идемпотентными.

---

## 📄 Swagger

### При запуске приложения можно просмотреть документацию:
* Swagger UI: `http://localhost:8080/swagger/index.html`
* JSON: `http://localhost:8080/swagger/doc.json`

---

## ✅ Итог

* #### Воспроизводимая тестовая среда.
* #### Единый источник данных через локальную БД.
* #### Интеграция с HTTP/gRPC.
* #### Swagger UI.
