package services

import (
	"soaProject/internal/db/entities"
	"soaProject/internal/db/models"
	"sort"
	"time"

	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"

	"gorm.io/gorm"
)

type StatementService struct {
	DB        *gorm.DB
	SecretKey string
	Salt      int
}

func NewStatementService(db *gorm.DB, vp *viper.Viper) *StatementService {
	return &StatementService{
		DB:        db,
		SecretKey: vp.GetString("jwt.secret"),
		Salt:      vp.GetInt("hash.salt"),
	}
}

func (ss *StatementService) GetStatement(ctx *fiber.Ctx) error {
	// Extract query parameters
	accountID := ctx.Query("accountID")
	startStr := ctx.Query("start")
	endStr := ctx.Query("end")

	// Validate the required fields
	if accountID == "" || startStr == "" || endStr == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: accountID, start, or end",
		})
	}

	// Parse the start and end times
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid start time format",
		})
	}
	end, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid end time format",
		})
	}

	var clientID string
	if err := ss.DB.Table("accounts").Where("id = ?", accountID).Select("client_id").Scan(&clientID).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve client ID",
		})
	}

	var clientName string
	if err := ss.DB.Table("clients").Where("id = ?", clientID).Select("name").Scan(&clientName).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve client name",
		})
	}

	var activity []models.ActivityInfo

	// Query the transactions from the database
	var transactions []entities.Transaction
	if err := ss.DB.Where("sender_id = ? OR receiver_id = ? AND created_at BETWEEN ? AND ?", accountID, accountID, start, end).Find(&transactions).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	for _, tx := range transactions {
		activity = append(activity, models.ActivityInfo{
			TxID:      tx.ID,
			From:      tx.SenderID,
			To:        tx.ReceiverID,
			Amount:    tx.Amount,
			Type:      "transaction",
			Timestamp: tx.CreatedAt,
		})
	}

	// Query the payments from the database
	var payments []entities.Payment
	if err := ss.DB.Where("account_id = ? AND created_at BETWEEN ? AND ?", accountID, start, end).Find(&payments).Error; err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve payments",
		})
	}

	for _, tx := range payments {
		activity = append(activity, models.ActivityInfo{
			TxID:      tx.ID,
			From:      tx.AccountID,
			To:        tx.RefCode,
			Amount:    tx.Amount,
			Type:      "payment",
			Timestamp: tx.CreatedAt,
		})
	}

	// Sort the activity by timestamp
	sort.Slice(activity, func(i, j int) bool {
		return activity[i].Timestamp.Before(activity[j].Timestamp)
	})

	// Create the statement response
	statement := models.StatementInfo{
		ClientName: clientName,
		ClientID:   clientID,
		AccountID:  accountID,
		Start:      start,
		End:        end,
		Activity:   activity,
	}

	// Return the statement as a JSON response
	return ctx.Status(fiber.StatusOK).JSON(statement)
}
