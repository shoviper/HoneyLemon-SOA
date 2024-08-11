package main

import (
	"flag"

	"github.com/Nukie90/SOA-Project/server"
)

func main() {
	env := flag.String("env", "dev", "Environment")
	flag.Parse()

	go server.ESBServer(*env)
	go server.AccountServer(*env)
	go server.PaymentServer(*env)
	go server.ClientServer(*env)
	go server.StatementServer(*env)
	go server.TransactionServer(*env)
	select {}
}
