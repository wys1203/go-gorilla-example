.PHONY: run-postgres stop-postgres build run clean

run-postgres:
	docker run --name postgres-ui-test \
		-e POSTGRES_DB=ui_test \
		-e POSTGRES_USER=ui_test \
		-e POSTGRES_PASSWORD=mysecretpassword \
		-p 5432:5432 \
		-d postgres

stop-postgres:
	docker stop postgres-ui-test
	docker rm postgres-ui-test

build:
	@echo "Building Go application..."
	@go build -o main ./cmd

run:
	@echo "Running Go application..."
	@go run .

clean:
	@echo "Cleaning up..."
	@rm -f main

docker-build:
	@echo "Building Docker image..."
	@docker-compose build

docker-up:
	@echo "Starting application and PostgreSQL services with Docker Compose..."
	@docker-compose up -d

docker-down:
	@echo "Stopping application and PostgreSQL services with Docker Compose..."
	@docker-compose down

docker-logs:
	@echo "Displaying logs of application and PostgreSQL services..."
	@docker-compose logs -f
