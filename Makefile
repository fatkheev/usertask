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
