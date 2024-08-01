package api

import (
	"soaProject/api/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	clientService := services.NewClientService(db)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.Get("/clients", clientService.GetAllClients)
			v1.Post("/clients", clientService.RegisterClient)
		}
	}
}