package main

import (
	"github.com/Nukie90/SOA-Project/server"
)

func main() {
	server.ClientServer("env", "dev", "a string")
}