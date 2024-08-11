package statement

import (
	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupStatementRoute(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	ss := NewStatementService(db)

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			statement := v1.Group("/statements")
			{
				statement.Get("/", ss.GetStatement)
			}
		}
	}

}
