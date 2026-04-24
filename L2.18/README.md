# Calendar HTTP Service

Микросервис для управления календарем событий. Реализован на Go с соблюдением принципов чистой архитектуры.

## Особенности
- **Слоистая архитектура**: чёткое разделение на транспортный слой (HTTP), бизнес-логику (Service) и хранилище (Repository).
- **Конфигурация**: использование библиотеки `cleanenv` (поддержка YAML и ENV).
- **In-memory хранилище**: потокобезопасная реализация на `sync.RWMutex`.
- **Валидация**: кастомная валидация входных данных с группировкой ошибок.
- **Middleware**: логирование каждого запроса (метод, URL, время выполнения).
- **Docker**: поддержка контейнеризации (multi-stage build).

## Требования
- Go 1.23+
- Docker (опционально)

## Конфигурация
Сервис использует файл конфигурации и переменные окружения.
- `CONFIG_PATH` — путь к YAML файлу (по умолчанию `config/local.yaml`).

**Основные параметры (можно переопределить через ENV):**
- `ENV` — окружение (`local`, `dev`, `prod`).
- `ADDRESS` — адрес и порт сервера (default: `0.0.0.0:8080`).
- `READ_TIMEOUT`, `WRITE_TIMEOUT`, `IDLE_TIMEOUT` — таймауты сервера.

## Запуск проекта

### Локально
```bash
make run
```

### Через Docker
```bash
make docker-build
make docker-run
```

## API Методы

### События (POST)
Передаются в формате JSON.
- `POST /create_event` — создание события.
- `POST /update_event` — обновление (требует ID).
- `POST /delete_event` — удаление (требует ID).
Согласно стандартам REST для update_event и delete_event следует использовать PUT и DELETE, но в данном проекте
используется POST в соответствии с требованиями ТЗ.

Пример тела запроса для создания (`POST /create_event`):
```json
{
  "user_id": 1,
  "date": "2024-10-25",
  "event": "Собеседование"
}
```

Пример тела запроса для обновления (`POST /update_event`):
```json
{
  "id": "46c4620b-a846-404c-91a5-7ec2871181a5",
  "user_id": 1,
  "date": "2024-10-25",
  "event": "Собеседование"
}
```

Пример тела запроса для удаления (`POST /delete_event`):
```json
{
  "id": "46c4620b-a846-404c-91a5-7ec2871181a5"
}
```

### Выборка (GET)
Параметры передаются через Query String: `?user_id=1&date=2024-10-25`
- `GET /events_for_day`
- `GET /events_for_week`
- `GET /events_for_month`

## Тестирование
Запуск unit-тестов сервисного слоя с проверкой на race condition:
```bash
make test
```

## Структура проекта
- `cmd/server/main.go` — точка входа, инициализация зависимостей.
- `internal/config` - загрузка конфигурации
- `internal/domain` — доменные сущности и ошибки.
- `internal/service` — бизнес-логика.
- `internal/repository` — хранение данных.
- `internal/transport/http` — хендлеры, роутинг, DTO и middleware.

