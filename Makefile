migrate-file:
	migrate create -ext sql -dir migrations ${filename}

wire:
	go run -mod=mod github.com/google/wire/cmd/wire ./internal

test:
	go test -v ./...

run-local:
	go run main.go start --config=./config/local/config.yaml

build:
	go build -o bin/shortbit main.go
