package main

import (
	"userService/api/grpc"
	"userService/api/server"
	"userService/internal/repository"
)

func main() {
	userService, err := repository.NewService("mongodb://localhost:27017")
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
