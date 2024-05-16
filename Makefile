build:
	@go build -o bin/realtime cmd/api/main.go
	@go build -o bin/cli cmd/cli/main.go

run:
	@./bin/realtime

migrateup:
	atlas migrate apply -u "${DATABASE_URL}" --dir=file://internal/database/migrations

migratehash:
	atlas migrate hash --dir=file://internal/database/migrations

migratenew:
	atlas migrate new "${NAME}" --dir=file://internal/database/migrations