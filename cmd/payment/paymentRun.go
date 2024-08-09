package main

import (
	"soaProject/server"
)

func main() {
	server.PaymentServer("env", "dev", "a string")
}