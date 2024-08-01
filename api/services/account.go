package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"

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
	var account models.AccountInfo
	if err := ctx.BodyParser(&account); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newAccount := entities.Account{
		ClientID: account.ClientID,
		Type:     account.Type,
		Balance:  account.Balance,
		Pin:      account.Pin,
	}

	if err := cs.accountDB.Create(&newAccount).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(newAccount)
}
