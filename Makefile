ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: fmt lint test
	go build -o ./bin/orderservice cmd/orderservice/main.go

fmt:)

	go fmt ./...

test:
	go test ./...

lint:
	golangci-lint run

up:
	docker-compose -f docker/docker-compose.yml up -d

down:
	docker-compose -f docker/docker-compose.yml down

db:
	mysql -h 127.0.0.1 -u $(ORDER_DATABASE_USER) -p$(ORDER_DATABASE_PASSWORD) $(ORDER_DATABASE_NAME)

migrate:
	migrate -database "$(ORDER_DATABASE_DRIVER)://$(ORDER_DATABASE_USER):$(ORDER_DATABASE_PASSWORD)@tcp(localhost:3370)/$(ORDER_DATABASE_NAME)" -path ./migrations up