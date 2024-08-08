package tokenbreaker

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func BreakToken(c *fiber.Ctx) error {
	fmt.Println("Token Breaker")
	cookie := c.Cookies("esb_token")
	if cookie =="" {
		fmt.Println("Token not found")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	//split token
	splitToken := strings.Split(cookie, ".")
	
	//merge only last 3 parts
	mergeToken := strings.Join(splitToken[1:], ".")

	fmt.Println("Token: ", mergeToken)

	c.Cookie(&fiber.Cookie{
		Name: "backend_token",
		Value: mergeToken,
		SameSite: "None", // Allow cross-site cookies
		Secure:   false,  // Not using HTTPS on localhost
	})

	return c.Next()
}