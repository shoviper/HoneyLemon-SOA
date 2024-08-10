package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type TransactionService struct {
	transactionDB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{transactionDB: db}
}

func (cs *TransactionService) GetAllTransactions(ctx *fiber.Ctx) error {
	var transactions []entities.Transaction

	if err := cs.transactionDB.Find(&transactions).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})

	}

	var transactionInfo []models.TransactionInfo
	for _, transaction := range transactions {
		transactionInfo = append(transactionInfo, models.TransactionInfo{
			ID:         transaction.ID,
			SenderID:   transaction.SenderID,
			ReceiverID: transaction.ReceiverID,
			Amount:     transaction.Amount,
		})
	}

	return ctx.Status(200).JSON(transactionInfo)
}

func (cs *TransactionService) CreateTransaction(ctx *fiber.Ctx) error {
	var transaction models.TransactionInfo
	if err := ctx.BodyParser(&transaction); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	newTransaction := entities.Transaction{
		ID:         transaction.ID,
		SenderID:   transaction.SenderID,
		ReceiverID: transaction.ReceiverID,
		Amount:     transaction.Amount,
	}

	if err := cs.transactionDB.Create(&newTransaction).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(200).JSON(newTransaction)
}
