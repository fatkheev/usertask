ifneq (,$(wildcard .env))
    include .env
    export
endif

.PHONY: run docker-build docker-run docker-stop docker-clean migrate-up migrate-down db-reset

run:
	go run cmd/main.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up --build -d

docker-stop:
	docker-compose down

docker-clean: docker-stop
	docker-compose down --volumes --remove-orphans
	docker rmi usertask || true

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@localhost:$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

db-reset:
	docker exec -it $(DB_CONTAINER) psql -U $(DB_USER) -d $(DB_NAME) -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
