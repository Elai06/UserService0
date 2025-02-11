# Makefile для Go-проекта

# Название исполнимого файла
BINARY_NAME=taskService

# Цель по умолчанию (что будет выполняться при запуске 'make')
all: build, lint, protoRun

# Сборка проекта
build: main.go
	go build -o $(BINARY_NAME) main.go

# Запуск тестов
test:
	go test ./...

lint:
	golangci-lint run

protoRun:
	protoc --go_out=. --go-grpc_out=. proto/userService.proto

# Очистка сгенерированных файлов
clean:
	rm -f $(BINARY_NAME)

# Запуск приложения
run: build
	./$(BINARY_NAME)
