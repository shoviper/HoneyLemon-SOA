package client

import (
	"client/src"

	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupClientRoute(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	cs := NewClientService(db, vp)

	auth := services.NewJWTConfig(vp)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			client := v1.Group("/clients")
			{
				client.Get("/", cs.GetAllClients)
				client.Post("/register", cs.RegisterClient)
				client.Post("/login", cs.LoginClient)
				client.Delete("/", cs.DeleteClient)
				client.Get("/logout", cs.LogoutClient)
				client.Get("/:id", cs.GetClientByID)
			}

			authClient := v1.Group("/authclients")
			{
				authClient.Use(auth.JWTAuth())
				authClient.Get("/info", cs.GetClientInfo)
			}
		}
	}

}
