package cmd

import (
	"flag"
	"fmt"
	"log"
	"soaProject/api"
	"soaProject/api/services"
	"soaProject/internal/config"
	"soaProject/internal/db"

	"soaProject/api/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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


	ServiceServer := fiber.New()

	ServiceServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: "http://localhost:3000, http://localhost:4000",
	}))

	api.SetupRoutes(ServiceServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(ServiceServer)
	
	serverConfig := config.NewServerConfig(configDetail)

	serverAddress := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)

	go func ()  {
		if err := ServiceServer.Listen(serverAddress); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	ESBServer := fiber.New()

	ESBServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins: "http://localhost:3000, http://localhost:4000",
	}))

	middleware.ESBRoute(ESBServer)

	esbConfig := config.NewEsbServerConfig(configDetail)

	esbAddress := fmt.Sprintf("%s:%d", esbConfig.Host, esbConfig.Port)

	go func ()  {
		if err := ESBServer.Listen(esbAddress); err != nil {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	<-make(chan struct{})
	return nil

}