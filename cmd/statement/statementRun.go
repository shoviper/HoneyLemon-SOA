package main

import (
	"soaProject/server"
)

func main() {
	server.StatementServer("env", "dev", "a string")
}
