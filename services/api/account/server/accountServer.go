package server

import (
	_ "flag"
	"log"

	"account/internal/config"
	"account/internal/db"
	"account/src"
	acc "account/src/account"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AccountServer(env string) error {
	configDetail, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbConfig := db.NewConfig(configDetail)
	db, err := dbConfig.PostgresConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	accountServer := fiber.New()

	accountServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://127.0.0.1:4000",
	}))

	acc.SetupAccountRoute(accountServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(accountServer)

	if err := accountServer.Listen(":3002"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
