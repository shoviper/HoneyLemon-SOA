package main

import (
	"soaProject/server"
)

func main() {
	server.TransactionServer("env", "dev", "a string")
}