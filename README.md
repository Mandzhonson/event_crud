# Calendar API

HTTP-сервер для управления событиями календаря. Реализует CRUD операции над событиями с хранением данных в PostgreSQL.

## Стек

- **Go** — основной язык
- **Gin** — HTTP-фреймворк
- **PostgreSQL** — база данных
- **Docker / Docker Compose** — контейнеризация
- **golang-migrate** — миграции базы данных
- **gomock + testify** — unit-тестирование

## Структура проекта

```
├── cmd/
│   └── main.go               # Точка входа
├── internal/
│   ├── app/                  # Инициализация приложения
│   ├── apperr/               # Ошибки приложения
│   ├── config/               # Конфигурация из .env
│   ├── dto/                  # DTO для запросов
│   ├── handlers/             # HTTP-обработчики + тесты
│   ├── logger/               # Настройка логгера
│   ├── middleware/           # Middleware логирования
│   ├── models/               # Модели данных
│   ├── repository/           # Слой работы с БД + моки
│   ├── router/               # Регистрация роутов
│   ├── service/              # Бизнес-логика + тесты
│   └── storage/              # Подключение к PostgreSQL
├── migrations/               # SQL-миграции
├── docker-compose.yml
├── Dockerfile
├── Makefile
└── .env
```

## API

Базовый путь: `/api/v1`

| Метод    | Путь           | Описание                        |
|----------|----------------|---------------------------------|
| `POST`   | `/events`      | Создать событие                 |
| `PUT`    | `/events/:id`  | Обновить событие                |
| `DELETE` | `/events/:id`  | Удалить событие                 |
| `GET`    | `/events`      | Получить события за период      |

### POST /api/v1/events

Создать новое событие.

**Тело запроса:**
```json
{
  "user_id": 1,
  "event_date": "2024-01-01",
  "event": "Meeting"
}
```

**Успешный ответ** `201 Created`:
```json
{
  "result": {
    "event_id": 1,
    "user_id": 1,
    "event_date": "2024-01-01",
    "event": "Meeting"
  }
}
```

### PUT /api/v1/events/:id

Обновить существующее событие.

**Тело запроса:**
```json
{
  "event_date": "2024-02-01",
  "event": "Updated Meeting"
}
```

**Успешный ответ** `200 OK`:
```json
{
  "result": "event updated"
}
```

### DELETE /api/v1/events/:id

Удалить событие по ID.

**Успешный ответ** `204 No Content`

### GET /api/v1/events

Получить события за период.

**Query параметры:**

| Параметр     | Тип    | Описание                          |
|--------------|--------|-----------------------------------|
| `period`     | string | Период: `day`, `week`, `month`    |
| `user_id`    | int    | ID пользователя                   |
| `event_date` | string | Дата в формате `YYYY-MM-DD`       |

**Пример запроса:**
```
GET /api/v1/events?period=day&user_id=1&event_date=2024-01-01
```

**Успешный ответ** `200 OK`:
```json
[
  {
    "event_id": 1,
    "user_id": 1,
    "event_date": "2024-01-01T00:00:00Z",
    "event": "Meeting"
  }
]
```

### Коды ошибок

| Код  | Описание                        |
|------|---------------------------------|
| 400  | Некорректные параметры запроса  |
| 404  | Событие не найдено              |
| 500  | Внутренняя ошибка сервера       |

**Формат ошибки:**
```json
{
  "error": "описание ошибки"
}
```

## Запуск

### Переменные окружения

Создай файл `.env` в корне проекта:

```env
SRV_HOST=0.0.0.0
SRV_PORT=8080

POSTGRES_DB=crud_app
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=postgresql
POSTGRES_PORT=5432

MIGRATE_HOST=localhost
```

### Docker Compose

```bash
# Запустить приложение и базу данных
docker compose up -d

# Остановить
docker compose down

# Остановить и удалить volumes
docker compose down -v
```

### Миграции

```bash
# Применить миграции
make migrate-up

# Откатить последнюю миграцию
make migrate-down
```

### Локальный запуск (без Docker)

```bash
go run cmd/main.go
```

## Тесты

```bash
# Запустить тесты
make test

# Покрытие в терминале
make coverage

# Покрытие в HTML
make coverage-html
```
