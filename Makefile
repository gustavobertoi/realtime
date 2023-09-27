build:
	@go build -o bin/realtime cmd/main.go

run:
	@./bin/realtime

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f