package account

import (
	"fmt"

	"account/internal/db/entities"
	"account/internal/db/models"
	local "account/internal/local"

	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"gorm.io/gorm"
)

type AccountService struct {
	accountDB *gorm.DB
	Salt      int
}

func NewAccountService(db *gorm.DB, vp *viper.Viper) *AccountService {
	return &AccountService{
		accountDB: db,
		Salt:      vp.GetInt("hash.salt"),
	}
}

func (as *AccountService) GetAllAccounts(ctx *fiber.Ctx) error {
	fmt.Println("Get all accounts")
	var accounts []entities.Account

	if err := as.accountDB.Find(&accounts).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get accounts",
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

func (as *AccountService) CreateAccount(ctx *fiber.Ctx) error {
	var account models.CreateAccount
	if err := ctx.BodyParser(&account); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	userID := ctx.Locals("userID")

	hasher := local.NewLocalConfig(as.Salt)

	//hash the pin
	hashedPin, err := hasher.HashPassword(account.Pin)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	account.Balance = 0

	newAccount := entities.Account{
		ClientID: userID.(string),
		Type:     account.Type,
		Balance:  account.Balance,
		Pin:      hashedPin,
	}

	if err := as.accountDB.Create(&newAccount).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "Account created successfully",
		"account": newAccount,
	})
}

func (as *AccountService) GetAccount(ctx *fiber.Ctx) error {
	fmt.Println("Get account by ID")
	var account models.AccountVerify
	account.ID = ctx.Params("id")
	fmt.Println(account.ID)
	if err := ctx.BodyParser(&account); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	var accountInfo entities.Account
	if err := as.accountDB.Where("id = ?", account.ID).First(&accountInfo).Error; err != nil {
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

func (as *AccountService) GetAllClientAccounts(ctx *fiber.Ctx) error {
	fmt.Println("Get all client's accounts")
	userID := ctx.Locals("userID")

	var accounts []entities.Account
	if err := as.accountDB.Where("client_id = ?", userID).Find(&accounts).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to get accounts",
		})
	}

	var accountsInfo []models.AccountBalance
	for _, account := range accounts {
		accountsInfo = append(accountsInfo, models.AccountBalance{
			ID:      account.ID,
			Type:    account.Type,
			Balance: account.Balance,
		})
	}

	return ctx.Status(200).JSON(accountsInfo)
}

func (as *AccountService) ChangePin(ctx *fiber.Ctx) error {
	var pinChange struct {
		AccountID string `json:"accountID"`
		OldPin    string `json:"oldPin"`
		NewPin    string `json:"newPin"`
	}

	// Parse the request body
	if err := ctx.BodyParser(&pinChange); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	// Retrieve the userID from the context
	userID := ctx.Locals("userID")

	var accountInfo entities.Account
	if err := as.accountDB.Where("id = ? AND client_id = ?", pinChange.AccountID, userID).First(&accountInfo).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to retrieve account information",
		})
	}

	// Verify the current PIN
	if !local.CheckPasswordHash(pinChange.OldPin, accountInfo.Pin) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid current PIN",
		})
	}

	// Hash the new PIN
	hasher := local.NewLocalConfig(as.Salt)
	hashedNewPin, err := hasher.HashPassword(pinChange.NewPin)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to hash new PIN",
		})
	}

	// Update the PIN in the database
	if err := as.accountDB.Model(&accountInfo).Update("pin", hashedNewPin).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to update PIN",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "PIN changed successfully",
	})
}

func (as *AccountService) DeleteAccount(ctx *fiber.Ctx) error {
	var delAccount struct {
		AccountID string `json:"accountID"`
		Pin       string `json:"pin"`
	}

	// Parse the request body
	if err := ctx.BodyParser(&delAccount); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Invalid request body",
		})
	}

	// Retrieve the userID from the context
	userID := ctx.Locals("userID")

	var accountInfo entities.Account
	if err := as.accountDB.Where("id = ? AND client_id = ?", delAccount.AccountID, userID).First(&accountInfo).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to retrieve account information",
		})
	}

	// Verify the current PIN
	if !local.CheckPasswordHash(delAccount.Pin, accountInfo.Pin) {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid current PIN",
		})
	}

	// Delete the account
	if err := as.accountDB.Delete(&accountInfo).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   err.Error(),
			"message": "Failed to delete account",
		})
	}

	return ctx.Status(200).JSON(fiber.Map{
		"message": "Account deleted successfully",
	})
}