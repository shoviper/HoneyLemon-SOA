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
			client := v1.Group("/clients")
			{
				client.Get("/", clientService.GetAllClients)
				client.Post("/register", clientService.RegisterClient)
				client.Post("/login", clientService.LoginClient)
				client.Delete("/", clientService.DeleteClient)
			}
			account := v1.Group("/accounts")
			{
				account.Use(services.JWTAuth())
				account.Get("/", accountService.GetAllAccounts)
				account.Post("/", accountService.CreateAccount)
				account.Get("/:id", accountService.GetAccount)
			}
			transaction := v1.Group("/transactions")
			{
				transaction.Get("/", WrapHandler(transactionService.GetAllTransactions))
				transaction.Post("/", WrapHandler(transactionService.CreateTransaction))
			}
		}
	}
}
