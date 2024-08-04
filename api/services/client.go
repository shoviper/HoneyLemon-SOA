package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"
	local "soaProject/internal/local"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	viper "github.com/spf13/viper"

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
	var client models.RegisterClient
	if err := ctx.BodyParser(&client); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Invalid request body",
		})
	}

	// Convert string to time
	birthDay, err := time.Parse("2006-01-02", client.BirthDate)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
			"Birthday": client.BirthDate,
		})
	}

	// Hash password
	hashPassword, err := local.HashPassword(client.Password)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}


	newClient := entities.Client{
		ID:        client.ID,
		Name:      client.Name,
		Address:   client.Address,
		BirthDate: birthDay,
		Password:  hashPassword,
	}

	if err := cs.clientDB.Create(&newClient).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(newClient)
}

func (cs *ClientService) LoginClient(ctx *fiber.Ctx) error {
	var clientReq models.LoginClient
	if err := ctx.BodyParser(&clientReq); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Invalid request body",
		})
	}

	var clientDB entities.Client
	if err := cs.clientDB.Where("id = ?", clientReq.ID).First(&clientDB).Error; err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Client not found",
		})
	}

	if !local.CheckPasswordHash(clientReq.Password, clientDB.Password) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid password",
		})
	}

	claims := jwt.MapClaims{
		"id": clientDB.ID,
		"exp": time.Now().Add(time.Hour * 24 * 365).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(viper.GetString("jwt.secret")))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(fiber.Map{

		"token": tokenString,
	})
}