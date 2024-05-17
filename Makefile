build:
	@go build -o bin/realtime cmd/api/main.go
	@go build -o bin/cli cmd/cli/main.go

run:
	@go build -o bin/realtime cmd/api/main.go
	@./bin/realtime

artillery:
	artillery run tests/artillery.yaml --output tests/artifacts.json
