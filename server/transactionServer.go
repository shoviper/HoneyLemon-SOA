package server

import (
	"flag"
	"fmt"
	"log"
	"soaProject/api/services"
	"soaProject/internal/config"
	"soaProject/internal/db"

	trans "soaProject/api/services/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func TransactionServer(name, value, usage string) error {
	envFlag := flag.String(name, value, usage)

	flag.Parse()

	configDetail, err := config.LoadConfig(*envFlag)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbConfig := db.NewConfig(configDetail)
	db, err := dbConfig.PostgresConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	transactionServer := fiber.New()

	transactionServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:4000",
	}))

	trans.SetupTransactionRoute(transactionServer, db, configDetail)
	transactionConfig := config.NewTransactionServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(transactionServer)

	transactionAddress := fmt.Sprintf("%s:%d", transactionConfig.Host, transactionConfig.Port)

	if err := transactionServer.Listen(transactionAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}