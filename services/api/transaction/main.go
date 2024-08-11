package main

import (
	"flag"

	"transaction/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.TransactionServer(*env)
}
