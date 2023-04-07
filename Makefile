up:
	docker-compose up -d --build

down:
	docker-compose down

lint:
	golangci-lint run

test:
	go test -race ./...
