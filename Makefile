SHELL := /bin/bash

build-app:
	@go build -o bin/go-ekyc cmd/app/main.go

run-app: build-app
	@./bin/go-ekyc

build-migrate:
	@go build -o bin/migrate cmd/migration/migration.go

build-worker:
	@go build -o bin/go-ekyc-worker cmd/worker/worker.go

worker: create-worker
	@./bin/go-ekyc-worker

build-cronjob:
	@go build -o bin/go-ekyc-cronjob cmd/cronjob/main.go

cronjob: create-cronjob
	@./bin/go-ekyc-cronjob

test-unit:
	@go test ./handler ./service ./cronjob ./worker

lint:
	@gofmt -l -s .

build:
	@ENV_FILE=.docker-compose.env docker compose build
	# @set -o allexport
	# @source .docker-compose.env
	# @set +o allexport
	# @docker compose exec database sh -c 'until pg_isready -U postgres; do sleep 1; done'
	# @sleep 30
	# # @docker compose logs -f # Attach logs for all services
	# # @make create-migrate # Build the migrate binary
	# # @./bin/migrate -v 5 -f # Run the migration

test-integration:
	@go test -v ./test/integration

load-test-face:
	@artillery run test/load/load_test_face_match.yml --output testdata/result_face_match.json

load-test-ocr:
	@artillery run test/load/load_test_ocr.yml --output testdata/result_ocr.json

test-coverage:
	@go test -coverprofile=testdata/coverage.out ./... && go tool cover -func=testdata/coverage.out
