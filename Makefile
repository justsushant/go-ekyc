build:
	@go build -o bin/go-ekyc cmd/main.go

run: build
	@./bin/go-ekyc