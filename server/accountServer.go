package server

import (
	_ "flag"
	"fmt"
	"log"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"

	acc "github.com/Nukie90/SOA-Project/api/services/account"

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
		AllowOrigins:     "http://localhost:4000",
	}))

	acc.SetupAccountRoute(accountServer, db, configDetail)
	accountConfig := config.NewAccServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(accountServer)

	accountAddress := fmt.Sprintf("%s:%d", accountConfig.Host, accountConfig.Port)

	if err := accountServer.Listen(accountAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
