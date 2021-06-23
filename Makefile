ifneq (,$(wildcard ./.env))
    include .env
    export
endif

build: modules fmt lint test proto
	go build -o ./bin/orderservice cmd/main.go
	docker-compose -p orderservice -f docker/docker-compose.yml build

fmt:
	go fmt ./...

test:
	go test ./...

lint:
	golangci-lint run

up:
	docker-compose -p orderservice -f docker/docker-compose.yml up -d

down:
	docker-compose -p orderservice -f docker/docker-compose.yml down

db:
	mysql -h 127.0.0.1 -u $(ORDER_DATABASE_USER) -p$(ORDER_DATABASE_PASSWORD) $(ORDER_DATABASE_NAME)

migrate_up:
	migrate -database "$(ORDER_DATABASE_DRIVER)://$(ORDER_DATABASE_USER):$(ORDER_DATABASE_PASSWORD)@tcp(localhost:3370)/$(ORDER_DATABASE_NAME)" -path ./migrations up

migrate_down:
	migrate -database "$(ORDER_DATABASE_DRIVER)://$(ORDER_DATABASE_USER):$(ORDER_DATABASE_PASSWORD)@tcp(localhost:3370)/$(ORDER_DATABASE_NAME)" -path ./migrations down -all

logs:
	docker-compose -p orderservice -f docker/docker-compose.yml logs

modules:
	go mod tidy

proto:
	cd ./api/orderservice && protoc -I/usr/local/include -I. \
		 -I$$GOPATH/src \
		 -I. \
		 -I$$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
		 --swagger_out=logtostderr=true:. \
		 --grpc-gateway_out=logtostderr=true:. \
		 --go_out=plugins=grpc:. ./orderservice.proto
	cd ./api/storeservice && protoc -I/usr/local/include -I. \
		 -I$$GOPATH/src \
		 -I. \
		 -I$$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.16.0/third_party/googleapis \
		 --swagger_out=logtostderr=true:. \
		 --grpc-gateway_out=logtostderr=true:. \
		 --go_out=plugins=grpc:. ./storeservice.proto