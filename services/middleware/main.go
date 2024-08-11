package main

import (
	"flag"

	"middleware/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.ESBServer(*env)
}
