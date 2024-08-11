package server

import (
	_ "flag"
	"log"

	"payment/internal/config"
	"payment/internal/db"
	"payment/src"
	pay "payment/src/payment"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func PaymentServer(env string) error {
	configDetail, err := config.LoadConfig(env)
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
		AllowOrigins:     "http://127.0.0.1:4000",
	}))

	pay.SetupPaymentRoute(paymentServer, db, configDetail)

	JwtSetup := services.NewJWTConfig(configDetail)
	JwtSetup.JWT_Setup(paymentServer)

	if err := paymentServer.Listen(":3004"); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}

	return nil
}
