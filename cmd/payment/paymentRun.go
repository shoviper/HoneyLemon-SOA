package main

import (
	"github.com/Nukie90/SOA-Project/server"
)

func main() {
	server.PaymentServer("env", "dev", "a string")
}