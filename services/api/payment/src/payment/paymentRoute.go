package payment

import (
	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupPaymentRoute(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	ps := NewPaymentService(db)

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			payment := v1.Group("/payments")
			{
				payment.Post("/getAll", WrapHandler(ps.GetAllPayments))
				payment.Post("/getByID", WrapHandler(ps.GetPaymentByID))
				payment.Post("/getByAccountID", WrapHandler(ps.GetPaymentsByAccountID))
				payment.Post("/create", WrapHandler(ps.CreatePayment))
			}
		}
	}

}
