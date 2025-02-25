.PHONY: run stop clean docker-build docker-run docker-stop docker-clean

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

# Применить миграции
migrate-up:
	migrate -path migrations -database "postgres://usertask_user:usertask_password@localhost:5432/usertask_db?sslmode=disable" up

# Откатить миграции
migrate-down:
	migrate -path migrations -database "postgres://usertask_user:usertask_password@localhost:5432/usertask_db?sslmode=disable" down
	
# Очистить базу (полное удаление всех данных)
db-reset:
	docker exec -it usertask_db psql -U usertask_user -d usertask_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"