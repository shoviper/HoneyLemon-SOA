package main

import (
	"flag"

	"statement/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.StatementServer(*env)
}
