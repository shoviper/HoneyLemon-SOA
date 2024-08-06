package traffic

import (
	"bytes"
	"io/ioutil"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CheckRegisterClient(ctx *fiber.Ctx) error {
	requestBody := ctx.Body()

	// Make the request to the second service
	resp, err := http.Post("http://localhost:3000/api/v1/clients/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()

	// Read the response body from the second service
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).Send(body)
}