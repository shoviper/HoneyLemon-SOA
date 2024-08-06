package middleware

import (
	"fmt"

	tf "soaProject/api/middleware/traffic"

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
	}
}