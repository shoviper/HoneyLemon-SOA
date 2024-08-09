package main

import (
	"fmt"
	"soaProject/server"
)

func main() {
	fmt.Println("Hello, World!")
	// cmd.Server("env", "dev", "a string")
	server.ClientServer("env", "dev", "a string")
	// cmd.ESBServer("env", "dev", "a string")
}