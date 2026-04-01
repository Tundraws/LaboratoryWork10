# Лабораторная работа №10

## Студент

- ФИО: Мельникова А. С.
- Группа: 220032-11
- Вариант: 3

## Выполненные задания

### Средняя сложность

- Задание 3: реализована валидация входных данных в Go через Gin binding.
- Задание 5: реализована передача сложных JSON-структур между Go и Python сервисами.
- Задание 7: реализован graceful shutdown в Go и Python сервисах.

### Повышенная сложность

- Задание 3: добавлена JWT-аутентификация в Go-сервисе и проверка токенов в Python.
- Задание 5: оба сервиса развёрнуты через Docker Compose в общей сети.

## Структура проекта

- `go-service/` - Gin-сервис с JWT, валидацией и graceful shutdown.
- `python-service/` - FastAPI-сервис, который валидирует JWT и передаёт JSON в Go-сервис.
- `docker-compose.yml` - совместный запуск обоих сервисов.
- `PROMPT_LOG.md` - лог работы с ИИ.

## Локальный запуск

### Go-сервис

```powershell
cd go-service
go mod tidy
go test ./... -coverprofile=coverage.out
go run ./cmd/server
```

Go-сервис будет доступен по адресу `http://localhost:8080`.

### Python-сервис

Требуется установленный Python 3.12+.

```powershell
cd python-service
python -m venv .venv
.venv\Scripts\Activate.ps1
pip install -r requirements.txt
pytest
uvicorn app.main:app --host 0.0.0.0 --port 8000
```

Python-сервис будет доступен по адресу `http://localhost:8000`.

## Запуск через Docker Compose

Требуется установленный Docker Desktop.

```powershell
docker compose up --build
```

После запуска:

- Go-сервис: `http://localhost:8080`
- Python-сервис: `http://localhost:8000`

## Примеры запросов

### Получение JWT в Go

```http
POST /auth/token
Content-Type: application/json

{
  "username": "student",
  "password": "securepass123"
}
```

### Передача JSON через Python в Go

```http
POST /api/forward
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "request_id": "3fa85f64-5717-4562-b3fc-2c963f66afa6",
  "customer": "Ivan Petrov",
  "address": {
    "city": "Moscow",
    "street": "Tverskaya 1",
    "zip_code": "123456"
  },
  "items": [
    {
      "name": "keyboard",
      "quantity": 2,
      "price": 1500.5
    }
  ],
  "metadata": {
    "priority": "high",
    "tags": ["study", "api"]
  }
}
```
