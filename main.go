package main

import (
	"os"

	"userService/api/grpc"
	"userService/api/server"
	"userService/env"
	"userService/internal/repository"
)

func main() {
	err := env.LoadEnv()
	if err != nil {
		panic(err)
	}

	userService, err := repository.NewService(os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	hh := server.NewHTTPHandler(userService)

	err = hh.StartServer()
	if err != nil {
		panic(err)
	}

	err = grpc.Listen(userService)
	if err != nil {
		panic(err)
	}
}
