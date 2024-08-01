package cmd

import (
	"flag"
	"fmt"
	"log"
	"soaProject/api"
	"soaProject/internal/db"
	"soaProject/internal/config"

	"github.com/gofiber/fiber/v2"
)

func Server(name, value, usage string) error{
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

	app := fiber.New()

	api.SetupRoutes(app, db)

	serverConfig := config.NewServerConfig(configDetail)

	serverAddress := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	log.Fatal(app.Listen(serverAddress))

	return nil

}