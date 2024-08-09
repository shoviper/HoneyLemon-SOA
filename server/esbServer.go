package server

import (
	"flag"
	"fmt"
	"log"
	"soaProject/api/services"
	"soaProject/internal/config"

	"soaProject/api/middleware"
	
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ESBServer(name, value, usage string) error {
	envFlag := flag.String(name, value, usage)

	flag.Parse()

	configDetail, err := config.LoadConfig(*envFlag)
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	ESBServer := fiber.New()

	ESBServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: "http://localhost:3000, http://localhost:4000, http://localhost:3002",
	}))

	middleware.ESBRoute(ESBServer)

	esbConfig := config.NewEsbServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(ESBServer)

	esbAddress := fmt.Sprintf("%s:%d", esbConfig.Host, esbConfig.Port)

	if err := ESBServer.Listen(esbAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
