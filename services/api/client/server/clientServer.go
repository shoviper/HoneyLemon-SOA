package server

import (
	_ "flag"
	"log"

	"client/internal/config"
	"client/internal/db"
	"client/src"
	client "client/src/client"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ClientServer(env string) error {
	configDetail, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	dbConfig := db.NewConfig(configDetail)
	db, err := dbConfig.PostgresConnection()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	ClientServer := fiber.New()

	ClientServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://127.0.0.1:4000",
	}))

	client.SetupClientRoute(ClientServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(ClientServer)

	if err := ClientServer.Listen(":3001"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
