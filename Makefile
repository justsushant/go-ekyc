build:
	@go build -o bin/go-ekyc cmd/app/main.go

run: build
	@./bin/go-ekyc

create-migrate:
	@go build -o bin/migrate cmd/migration/migration.go
