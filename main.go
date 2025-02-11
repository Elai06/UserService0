package main

import (
	"os"
	"userService/api/grpc"
	"userService/api/server"
	"userService/env"
	"userService/internal/repository"
)

func main() {
	envErr := env.LoadEnv()
	if envErr != nil {
		panic(envErr)
	}

	userService, err := repository.NewService(os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	hh := server.NewHttpHandler(userService)
	hhErr := hh.StartServer()
	if hhErr != nil {
		panic(hhErr)
	}

	grpcErr := grpc.StartRpc(userService)
	if grpcErr != nil {
		panic(grpcErr)
	}
}
