package api

import (
	"soaProject/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	clientService := services.NewClientService(db)
	accountService := services.NewAccountService(db)
	transactionService := services.NewTransactionService(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.Get("/clients", clientService.GetAllClients)
			v1.Post("/clients", clientService.RegisterClient)

			v1.Get("/accounts", accountService.GetAllAccounts)
			v1.Post("/accounts", accountService.CreateAccount)

			v1.Get("/transactions", transactionService.GetAllTransactions)
			v1.Post("/transactions", transactionService.CreateTransaction)
		}
	}
}
