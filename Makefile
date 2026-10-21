migrate-file:
	migrate create -ext sql -dir migrations ${filename}

wire:
	go run github.com/google/wire/cmd/wire ./internal

test:
	go test -v ./...

run-local:
	go run main.go

build:
	go build -o bin/shortbit main.go
