.PHONY: run test migrate fmt build

run:
	air

migrate:
	go run main.go migrate

test:
	go test ./...

fmt:
	go fmt ./...

build:
	go build -o app .