package main

import (
	"flag"

	"client/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.ClientServer(*env)
}
