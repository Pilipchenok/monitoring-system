.PHONY: build run test race vet clean docker-up docker-down

build:
	go build -o bin/collector ./cmd/collector
	go build -o bin/agent ./cmd/agent

run-collector:
	go run ./cmd/collector

run-agent:
	go run ./cmd/agent

test:
	go test ./... -v

race:
	go test -race ./... -v

vet:
	go vet ./...

clean:
	rm -rf bin/

docker-up:
	docker compose up --build -d

docker-down:
	docker compose down -v
