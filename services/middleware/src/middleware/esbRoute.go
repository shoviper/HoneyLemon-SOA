package middleware

import (
	"fmt"

	// tf "github.com/Nukie90/SOA-Project/api/middleware/traffic"

	"github.com/gofiber/fiber/v2"
)

// ESBRoute is a function to setup ESB middleware
func ESBRoute(app *fiber.App) {
	app.Use(func(c *fiber.Ctx) error {
		fmt.Println("ESB middleware")
		return c.Next()
	})

	esb := app.Group("/esb")
	{
		esb.Post("/register", CheckRegisterClient)
		esb.Post("/login", CheckLoginClient)
		esb.Get("/logout", DoLogout)

		client := esb.Group("/clients")
		{
			client.Get("/info", GetClientInfo)
		}

		account := esb.Group("/accounts")
		{
			account.Get("/getAll", GetAllAccounts)
			account.Post("/create", CreateAccount)
			account.Get("/getByID/:id", GetAccount)
			account.Get("/clientAcc", GetAllClientAccounts)
		}

		transaction := esb.Group("/transactions")
		{
			transaction.Get("/getAll", GetAllTransactions)
			transaction.Get("/getByID", GetTransactionByID)
			transaction.Get("/getByAccountID", GetTransactionsByAccountID)
			transaction.Post("/create", CreateTransaction)
		}
		payment := esb.Group("/payments")
		{
			payment.Get("/getAll", GetAllPayments)
			payment.Get("/getByID", GetPaymentByID)
			payment.Get("/getByAccountID", GetPaymentsByAccountID)
			payment.Post("/create", CreatePayment)
		}
		statement := esb.Group("/statements")
		{
			statement.Get("/", GetStatement)
		}
	}
}
