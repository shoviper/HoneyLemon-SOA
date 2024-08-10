package middleware

import (
	"fmt"

	tf "github.com/Nukie90/SOA-Project/api/middleware/traffic"

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
		esb.Get("/logout", tf.DoLogout)

		account := esb.Group("/accounts")
		{
			account.Get("/getAll", tf.GetAllAccounts)
			account.Post("/create", tf.CreateAccount)
			account.Get("/getByID/:id", tf.GetAccount)
			account.Get("/clientAcc", tf.GetAllClientAccounts)
			account.Patch("/changePin", tf.ChangePin)
			account.Delete("/deleteAcc", tf.DeleteAccount)
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
		statement := esb.Group("/statements")
		{
			statement.Get("/", tf.GetStatement)
		}
	}
}
