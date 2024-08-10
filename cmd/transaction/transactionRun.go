package main

import (
	"github.com/Nukie90/SOA-Project/server"
)

func main() {
	server.TransactionServer("env", "dev", "a string")
}