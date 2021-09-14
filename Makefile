install:
	go get ./...

build:
	go build -o bin/main main.go

run:
	go run main.go

local:: build
	GOOS=linux GOARCH=amd64 go build -o bin/main main.go

api: build
	sam local start-api