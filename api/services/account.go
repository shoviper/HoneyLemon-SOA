package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"
	local "soaProject/internal/local"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type AccountService struct {
	accountDB *gorm.DB
}

func NewAccountService(db *gorm.DB) *AccountService {
	return &AccountService{accountDB: db}
}

func (cs *AccountService) GetAllAccounts(ctx *fiber.Ctx) error {
	var accounts []entities.Account

	if err := cs.accountDB.Find(&accounts).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}

	var accountsInfo []models.AccountInfo
	for _, account := range accounts {
		accountsInfo = append(accountsInfo, models.AccountInfo{
			ID:       account.ID,
			ClientID: account.ClientID,
			Type:     account.Type,
			Balance:  account.Balance,
			Pin:      account.Pin,
		})
	}

	return ctx.Status(200).JSON(accountsInfo)
}

func (cs *AccountService) CreateAccount(ctx *fiber.Ctx) error {
	var account models.CreateAccount
	if err := ctx.BodyParser(&account); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Invalid request body",
		})
	}

	userID := ctx.Locals("userID")

	//hash the pin
	hashedPin, err := local.HashPassword(account.Pin)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newAccount := entities.Account{
		ClientID: userID.(string),
		Type:     account.Type,
		Balance:  account.Balance,
		Pin:      hashedPin,
	}

	if err := cs.accountDB.Create(&newAccount).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Account created successfully",
		"account": newAccount,
	})
}

func (cs *AccountService) GetAccount(ctx *fiber.Ctx) error {
	var account models.AccountVerify
	account.ID = ctx.Params("id")
	if err := ctx.BodyParser(&account); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
			"message": "Invalid request body",
		})
	}

	var accountInfo entities.Account
	if err := cs.accountDB.Where("id = ?", account.ID).First(&accountInfo).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if !local.CheckPasswordHash(account.Pin, accountInfo.Pin) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid pin",
		})
	}

	if accountInfo.ClientID != ctx.Locals("userID") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "This user didn't own this account",
		})
	}

	accountBalance := models.AccountBalance{
		ID:      accountInfo.ID,
		Balance: accountInfo.Balance,
		Type:    accountInfo.Type,
	}

	return ctx.Status(200).JSON(accountBalance)
}