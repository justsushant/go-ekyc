build:
	@go build -o bin/go-ekyc cmd/app/main.go

run: build
	@./bin/go-ekyc

create-migrate:
	@go build -o bin/migrate cmd/migration/migration.go

create-worker:
	@go build -o bin/go-ekyc-worker cmd/worker/worker.go

worker: create-worker
	@./bin/go-ekyc-worker

create-cronjob:
	@go build -o bin/go-ekyc-cronjob cmd/cronjob/main.go

cronjob: create-cronjob
	@./bin/go-ekyc-cronjob

test:
	@go test ./...