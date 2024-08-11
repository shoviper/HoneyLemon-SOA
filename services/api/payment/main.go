package main

import (
	"flag"

	"payment/server"
)

func main() {
	env := flag.String("env", "pgsql", "Environment")
	flag.Parse()

	server.PaymentServer(*env)
}
