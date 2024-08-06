package cmd

import (
	"flag"
	"fmt"
	"log"
	"soaProject/api"
	"soaProject/api/services"
	"soaProject/internal/db"
	"soaProject/internal/config"

	"soaProject/api/middleware"

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

	app1 := fiber.New()

	api.SetupRoutes(app1, db)
	services.JWT_Setup(app1)
	
	serverConfig := config.NewServerConfig(configDetail)

	serverAddress := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	go func ()  {
		if err := app1.Listen(serverAddress); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	app2 := fiber.New()

	middleware.ESBRoute(app2)

	esbConfig := config.NewEsbServerConfig(configDetail)

	esbAddress := fmt.Sprintf("%s:%d", esbConfig.Host, esbConfig.Port)

	go func ()  {
		if err := app2.Listen(esbAddress); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-make(chan struct{})
	return nil

}