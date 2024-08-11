package server

import (
	_ "flag"
	"log"

	"middleware/internal/config"
	services "middleware/src"
	middleware "middleware/src/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ESBServer(env string) error {

	configDetail, err := config.LoadConfig(env)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ESBServer := fiber.New()

	ESBServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://client-services:3001, http://account-services:3002, http://transaction-services:3003, http://payment-services:3004, http://statement-services:3005, http://127.0.0.1:5000",
	}))

	middleware.ESBRoute(ESBServer)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(ESBServer)

	if err := ESBServer.Listen(":4000"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
