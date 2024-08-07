package services

import (
	"strings"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	viper "github.com/spf13/viper"
)

type JWT struct {
	Secret string
}

func NewJWTConfig(v *viper.Viper) *JWT {
	return &JWT{
		Secret: v.GetString("jwt.secret"),
	}
}

// JWT_Setup is a function to setup JWT middleware
func (j *JWT)JWT_Setup(app *fiber.App) {
	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(j.Secret)},
	}))
}

// JWTAuth is a function to authenticate JWT token
func (j *JWT)JWTAuth() func(*fiber.Ctx) error {
	return (func(c *fiber.Ctx) error {
		// read from cookie
		cookie := c.Cookies("token")
		if cookie == "" {
			return c.Status(401).JSON(fiber.Map{
				"message": "Unauthorized",
			})
		}

		accessToken := strings.Replace(cookie, "Bearer ", "", 1)

		token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(j.Secret), nil
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