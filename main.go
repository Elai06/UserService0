package main

import (
	"userService/api/server"
	"userService/internal/user"
)

func main() {
	user.ConnectToMongo()
	server.StartServer()
	server.StartRpc()
}
