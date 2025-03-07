# Usertask

Usertask — это REST API для управления пользователями и заданиями. Проект реализован на языке Go с использованием фреймворка Gin. API поддерживает регистрацию пользователей, начисление очков за выполнение заданий, установку рефералов, а также предоставляет лидерборд.

## Структура проекта

```
├── cmd
│   ├── main.go                         # Точка входа в приложение
├── internal
│   ├── auth                            # JWT-аутентификация
│   │   ├── jwt.go
│   ├── database                        # Подключение к базе данных
│   │   ├── db.go
│   ├── handlers                        # HTTP-обработчики
│   │   ├── math_handlers.go
│   │   ├── user_handlers.go
│   ├── middleware                      # Мидлвары
│   │   ├── auth_middleware.go
│   ├── models                          # Определение моделей данных
│   │   ├── requests.go
│   │   ├── response.go
│   │   ├── task.go
│   │   ├── user.go
│   ├── repository                      # Логика работы с БД
│   │   ├── user_repository.go
│   │   ├── user_repository_test.go
│   ├── service                         # Бизнес-логика
│   │   ├── math_service.go
│   │   ├── math_service_test.go
│   │   ├── user_service.go
│   │   ├── user_service_test.go
├── migrations                          # SQL-миграции для БД
│   ├── 000001_init_schema.up.sql
│   ├── 000001_init_schema.down.sql
├── docs                                # Документация API (Swagger)
├── .env                                # Файл конфигурации окружения
├── docker-compose.yml                  # Конфигурация Docker
├── Dockerfile                          # Docker-образ
├── Makefile                            # Автоматизация команд
├── go.mod                              # Модуль Go
├── go.sum                              # Зависимости Go
```

## Установка и запуск

### Зависимости

Перед запуском установите:
- Go 1.22+
- Docker
- PostgreSQL
- Make (для удобного запуска команд)

### Запуск в Docker

1. Соберите и запустите контейнеры:

```
make docker-build
make docker-run
```

2. Остановите контейнеры:

```
make docker-stop
```

## База данных

Используется PostgreSQL. Структура БД:

| Поле        | Тип                                         | Описание                   |
|------------|-------------------------------------------|----------------------------|
| id         | SERIAL PRIMARY KEY                        | Уникальный ID пользователя |
| username   | VARCHAR(255) UNIQUE NOT NULL             | Имя пользователя           |
| points     | INT DEFAULT 0                             | Очки пользователя          |
| referrer_id | INT REFERENCES users(id) ON DELETE SET NULL | ID реферала               |
| created_at | TIMESTAMP DEFAULT CURRENT_TIMESTAMP      | Дата создания              |

### Таблица tasks

| Поле         | Тип                                        | Описание              |
|-------------|------------------------------------------|----------------------|
| id         | SERIAL PRIMARY KEY                         | Уникальный ID задания |
| user_id    | INT REFERENCES users(id) ON DELETE CASCADE | ID пользователя       |
| task_type  | VARCHAR(255) NOT NULL                      | Тип задания           |
| points     | INT NOT NULL                               | Количество очков      |
| completed_at | TIMESTAMP DEFAULT CURRENT_TIMESTAMP       | Дата завершения       |


## Миграции

Установка golang-migrate

```
brew install golang-migrate
```

Применение миграций

```
make migrate-up
```

Откат миграций

```
make migrate-down
```

Сброс базы данных
```
make db-reset
```

## Swagger-документация

API документирован с помощью Swagger. Для генерации документации:

```
make swag-rebuild
```

Swagger доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

## API Методы

### 1. POST /users/create — создание пользователя
### Тело запроса:
```
{
  "username": "user123"
}
```
### Ответ:
```
{
  "user": {
    "id": 1,
    "username": "user123",
    "points": 0
  },
  "token": "<JWT_TOKEN>"
}
```
---

### 2. POST /users/token/refresh — обновление токена
### Тело запроса:
```
{
  "user_id": 1
}
```
### Ответ:
```
{
  "message": "new token generated",
  "token": "<JWT_TOKEN>"
}
```
---

### 3. GET /users/:id/status — получение информации о пользователе
### Ответ:
```
{
  "id": 1,
  "username": "user123",
  "points": 100
}
```
---

### 4. POST /users/:id/task/complete — завершение задания
### Тело запроса:
```
{
  "task_type": "task",
  "points": 50
}
```
### Ответ:
```
{
  "message": "task completed",
  "points_awarded": 50
}
```
---

### 5. POST /users/:id/referrer — установка реферала
### Тело запроса:
```
{
  "referrer_id": 2
}
```
### Ответ:
```
{
  "message": "referrer set successfully"
}
```
---

### 6. GET /users/:id/task/math — получение математической задачи
### Ответ:
```
{
  "operand1": 5,
  "operand2": 3,
  "operation": "+"
}
```
---

### 6. POST /users/:id/task/math/solve — отправка ответа на задачу
### Тело запроса:
```
{
  "answer": 8
}
```
### Ответ:
```
{
  "message": "correct answer!",
  "points_awarded": 50
}
```
---

### 7. GET /users/leaderboard — лидерборд
### Ответ:
```
[
  { "id": 1, "username": "Alice", "points": 500 },
  { "id": 2, "username": "Bob", "points": 450 }
]
```
---