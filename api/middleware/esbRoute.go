package middleware

import (
	"fmt"

	tf "soaProject/api/middleware/traffic"
	breaker "soaProject/api/middleware/tokenBreaker"

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
		esb.Post("/register", tf.CheckRegisterClient)
		esb.Post("/login", tf.CheckLoginClient)

		account := esb.Group("/accounts")
		{
			account.Get("/getAll", breaker.BreakToken ,tf.GetAllAccounts)
			account.Post("/create", breaker.BreakToken, tf.CreateAccount)
		}

		transaction := esb.Group("/transactions")
		{
			transaction.Get("/getAll", tf.GetAllTransactions)
			transaction.Get("/getByID", tf.GetTransactionByID)
			transaction.Get("/getByAccountID", tf.GetTransactionsByAccountID)
			transaction.Post("/create", tf.CreateTransaction)
		}
		payment := esb.Group("/payments")
		{
			payment.Get("/getAll", tf.GetAllPayments)
			payment.Get("/getByID", tf.GetPaymentByID)
			payment.Get("/getByAccountID", tf.GetPaymentsByAccountID)
			payment.Post("/create", tf.CreatePayment)
		}
	}
}
