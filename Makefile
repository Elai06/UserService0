BINARY_NAME=taskService

all: build, lint, lint-fix, protoRun

build: main.go
	go build -o $(BINARY_NAME) main.go

test:
	go test ./...

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix

protoRun:
	protoc --go_out=. --go-grpc_out=. proto/userService.proto

clean:
	rm -f $(BINARY_NAME)

run: build
	./$(BINARY_NAME)
