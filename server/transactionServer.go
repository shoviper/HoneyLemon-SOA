package server

import (
	"flag"
	"fmt"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"
	"log"

	trans "github.com/Nukie90/SOA-Project/api/services/transaction"

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
		AllowOrigins:     "http://localhost:4000, http://localhost:5000",
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
