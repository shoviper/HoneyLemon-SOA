package services

import (
	"strings"


	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	viper "github.com/spf13/viper"
)

// JWT_Setup is a function to setup JWT middleware
func JWT_Setup(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(viper.GetString("jwt.secret"))},
	}))
}

// JWTAuth is a function to authenticate JWT token
func JWTAuth() func(*fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		accessToken := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		if accessToken == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("jwt.secret")), nil
		})
		if err != nil {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		c.Locals("userID", claims["userID"])
		return c.Next()
	})
}
