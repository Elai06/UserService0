package main

import (
	"fmt"
	"os"
	"time"

	"userService/api/grpc"
	"userService/api/server"
	"userService/env"
	"userService/internal/repository"
)

const writeTimout = "WRITE_TIMEOUT"
const readTimout = "READ_TIMEOUT"
const mongoUrl = "MONGO_URL"

type config struct {
	writeTimeout time.Duration
	readTimeout  time.Duration
	mongoUrl     string
}

func main() {
	envConfig, err := initConfig()
	if err != nil {
		panic(err)
	}

	userService, err := repository.NewService(envConfig.mongoUrl, envConfig.writeTimeout)
	if err != nil {
		panic(err)
	}

	hh := server.NewHTTPHandler(userService)
	err = hh.StartServer(envConfig.writeTimeout, envConfig.readTimeout)
	if err != nil {
		panic(err)
	}

	err = grpc.Listen(userService)
	if err != nil {
		panic(err)
	}
}

func initConfig() (*config, error) {
	err := env.LoadEnv()
	if err != nil {
		panic(err)
	}

	writeTimeout, errEnv := env.GetTimeDuration(writeTimout)
	if errEnv != nil {
		return nil, fmt.Errorf("error while getting env vars: %w", errEnv)
	}

	readTimeout, errEnv := env.GetTimeDuration(readTimout)
	if errEnv != nil {
		return nil, fmt.Errorf("error while getting env vars: %w", errEnv)
	}

	return &config{
		writeTimeout: writeTimeout,
		readTimeout:  readTimeout,
		mongoUrl:     os.Getenv(mongoUrl),
	}, nil
}
