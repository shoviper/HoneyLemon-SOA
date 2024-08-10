package server

import (
	"flag"
	"fmt"
	"github.com/Nukie90/SOA-Project/api/services"
	"github.com/Nukie90/SOA-Project/internal/config"
	"github.com/Nukie90/SOA-Project/internal/db"
	"log"

	pay "github.com/Nukie90/SOA-Project/api/services/payment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func PaymentServer(name, value, usage string) error {
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

	paymentServer := fiber.New()

	paymentServer.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:4000, http://localhost:5000",
	}))

	pay.SetupPaymentRoute(paymentServer, db, configDetail)
	clientConfig := config.NewPaymentServerConfig(configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(paymentServer)

	paymentAddress := fmt.Sprintf("%s:%d", clientConfig.Host, clientConfig.Port)

	if err := paymentServer.Listen(paymentAddress); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
