package server

import (
	"flag"
	"fmt"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"
	"log"

	acc "github.com/Nukie90/SOA-Project/api/services/account"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func AccountServer(name, value, usage string) error {
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

	accountServer := fiber.New()

	accountServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://127.0.0.1:4000, http://127.0.0.1:5000",
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
