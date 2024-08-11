package server

import (
	_ "flag"
	"log"

	"transaction/internal/config"
	"transaction/internal/db"
	"transaction/src"
	trans "transaction/src/transaction"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func TransactionServer(env string) error {
	configDetail, err := config.LoadConfig(env)
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
		AllowOrigins:     "http://127.0.0.1:4000",
	}))

	trans.SetupTransactionRoute(transactionServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(transactionServer)

	if err := transactionServer.Listen(":3003"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
