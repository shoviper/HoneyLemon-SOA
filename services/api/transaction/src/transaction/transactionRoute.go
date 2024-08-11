package transaction

import (
	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupTransactionRoute(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	ts := NewTransactionService(db)

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			transaction := v1.Group("/transactions")
			{
				transaction.Post("/getAll", WrapHandler(ts.GetAllTransactions))
				transaction.Post("/getByID", WrapHandler(ts.GetTransactionByID))
				transaction.Post("/getByAccountID", WrapHandler(ts.GetTransactionsByAccountID))
				transaction.Post("/create", WrapHandler(ts.CreateTransaction))
			}
		}
	}

}
