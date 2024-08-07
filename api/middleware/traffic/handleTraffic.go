package traffic

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"soaProject/internal/db/models"

	"github.com/gofiber/fiber/v2"
)

func CheckRegisterClient(ctx *fiber.Ctx) error {
	requestBody := ctx.Body()

	//check if the request body is empty
	if len(requestBody) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is empty")
	}

	//check if the request body is valid JSON
	if !json.Valid(requestBody) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is not valid JSON")
	}

	//check input fields
	var client models.RegisterClient
	err := json.Unmarshal(requestBody, &client)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if client.Name == "" || client.Address == "" || client.BirthDate == "" || client.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	resp, err := http.Post("http://localhost:3000/api/v1/clients/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()

	// Read the response body from the second service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).JSON(fiber.Map{
		"message": "Response from second service",
		"body":    string(body),
	})
}