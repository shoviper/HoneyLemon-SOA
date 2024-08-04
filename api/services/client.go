package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ClientService struct {
	clientDB *gorm.DB
}

func NewClientService(db *gorm.DB) *ClientService {
	return &ClientService{clientDB: db}
}

func (cs *ClientService) GetAllClients(ctx *fiber.Ctx) error {
	var clients []entities.Client

	if err := cs.clientDB.Find(&clients).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}

	var clientsInfo []models.ClientInfo
	for _, client := range clients {
		birthdayStr := client.BirthDate.Format("2006-01-02")
		clientsInfo = append(clientsInfo, models.ClientInfo{
			ID:        client.ID,
			Name:      client.Name,
			Address:   client.Address,
			BirthDate: birthdayStr,
		})
	}

	return ctx.Status(200).JSON(clientsInfo)
}

func (cs *ClientService) RegisterClient(ctx *fiber.Ctx) error {
	var client models.ClientInfo
	if err := ctx.BodyParser(&client); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Invalid request body",
		})
	}

	BirthDay, err := time.Parse("2006-01-02", client.BirthDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"Birthday": client.BirthDate,
		})
	}

	newClient := entities.Client{
		ID:        client.ID,
		Name:      client.Name,
		Address:   client.Address,
		BirthDate: BirthDay,
	}

	if err := cs.clientDB.Create(&newClient).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(newClient)
}