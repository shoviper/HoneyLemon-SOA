package main

import (
	"fmt"
	"soaProject/server"

	"sync"
)

func main() {
	fmt.Println("Server is running")
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		defer wg.Done()
		server.ESBServer("env", "dev", "a string")
	}()

	go func() {
		defer wg.Done()
		server.AccountServer("env", "dev", "a string")
	}()

	go func() {
		defer wg.Done()
		server.ClientServer("env", "dev", "a string")
	}()

	wg.Wait()
}