# usertask
Pet


## Миграции

Для запуска проекта требуется установка golang-migrate

```brew install golang-migrate```

Также я установил в сам проект

```go get -u github.com/golang-migrate/migrate/v4```

В корне проекта создал папку migrations

```mkdir -p migrations```

Сгенерировал новую миграцию

```migrate create -ext sql -dir migrations -seq init_schema```


## База данных

Использую для БД Posgres

Объяснение структуры БД:

- Таблица users: id, username, points (баллы), referrer_id (реферал).
- Таблица tasks: id, user_id, task_type, points, completed_at.