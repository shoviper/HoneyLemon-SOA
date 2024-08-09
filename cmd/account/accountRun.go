package main

import (
	"github.com/Nukie90/SOA-Project/server"
)

func main() {
	server.AccountServer("env", "dev", "a string")
}