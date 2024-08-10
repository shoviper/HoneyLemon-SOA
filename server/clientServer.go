package server

import (
	"flag"
	"fmt"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"
	"log"

	client "github.com/Nukie90/SOA-Project/api/services/client"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func ClientServer(name, value, usage string) error {
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

	ClientServer := fiber.New()

	ClientServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:4000, http://localhost:5000",
	}))

	client.SetupClientRoute(ClientServer, db, configDetail)
	clientConfig := config.NewClientServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(ClientServer)

	clientAddress := fmt.Sprintf("%s:%d", clientConfig.Host, clientConfig.Port)

	if err := ClientServer.Listen(clientAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
