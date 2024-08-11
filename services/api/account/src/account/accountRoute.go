package account

import (
	"account/src"

	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupAccountRoute(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	as := NewAccountService(db, vp)

	auth := services.NewJWTConfig(vp)

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			accounts := v1.Group("/accounts")
			{
				accounts.Use(auth.JWTAuth())
				accounts.Get("/", as.GetAllAccounts)
				accounts.Post("/", as.CreateAccount)
				accounts.Get("/clientAcc/:id", as.GetAccount)
				accounts.Get("/clientAcc", as.GetAllClientAccounts)
				accounts.Patch("/changePin", as.ChangePin)
				accounts.Delete("/deleteAcc", as.DeleteAccount)
			}
		}
	}

}
