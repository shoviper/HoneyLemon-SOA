package server

import (
	"flag"
	"fmt"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"
	"log"

	statement "github.com/Nukie90/SOA-Project/api/services/statement"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StatementServer(name, value, usage string) error {
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

	statementServer := fiber.New()

	statementServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:4000, http://localhost:5000",
	}))

	statement.SetupStatementRoute(statementServer, db, configDetail)
	statementConfig := config.NewStatementServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(statementServer)

	statementAddress := fmt.Sprintf("%s:%d", statementConfig.Host, statementConfig.Port)

	if err := statementServer.Listen(statementAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
