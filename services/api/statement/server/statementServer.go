package server

import (
	_ "flag"
	"log"
	"statement/internal/config"
	"statement/internal/db"
	"statement/src"

	statement "statement/src/statement"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StatementServer(env string) error {
	configDetail, err := config.LoadConfig(env)
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
		AllowOrigins:     "http://127.0.0.1:4000",
	}))

	statement.SetupStatementRoute(statementServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(statementServer)

	if err := statementServer.Listen(":3005"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
