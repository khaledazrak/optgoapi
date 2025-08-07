APP_NAME := go-api

build:
	go build -o app main.go utils.go

run:
	go run main.go

test:
	go test -v

lint:
	golangci-lint run

docker-build:
	docker build -t $(APP_NAME) .

docker-run:
	docker run -p 8000:8000 $(APP_NAME)

