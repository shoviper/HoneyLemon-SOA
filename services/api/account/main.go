package main

import (
	"flag"

	"account/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.AccountServer(*env)
}
