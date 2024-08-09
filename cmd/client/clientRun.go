package main

import (
	"soaProject/server"
)

func main() {
	server.ClientServer("env", "dev", "a string")
}