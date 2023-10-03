build:
	@go build -o bin/realtime cmd/api/main.go

run:
	@./bin/realtime

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f