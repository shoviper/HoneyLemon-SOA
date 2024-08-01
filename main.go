package main

import (
	"fmt"
	"soaProject/cmd"
)

func main() {
	fmt.Println("Hello, World!")
	cmd.Server("env", "dev", "a string")
}