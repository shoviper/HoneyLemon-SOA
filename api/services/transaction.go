package services

import (
	"encoding/xml"
	"net/http"
	"soaProject/internal/db/entities"
	"time"

	"gorm.io/gorm"
)

type CreateTransactionRequest struct {
	XMLName     xml.Name        `xml:"CreateTransactionRequest"`
	Transaction TransactionInfo `xml:"transaction"`
}

type CreateTransactionResponse struct {
	Transaction TransactionInfo `xml:"transaction"`
}

type GetAllTransactionsResponse struct {
	Transactions []TransactionInfo `xml:"transactions"`
}

type TransactionInfo struct {
	ID         string    `xml:"ID"`
	SenderID   string    `xml:"SenderID"`
	ReceiverID string    `xml:"ReceiverID"`
	Amount     float64   `xml:"Amount"`
	CreatedAt  time.Time `xml:"CreatedAt"`
}

type TransactionService struct {
	DB *gorm.DB
}

func NewTransactionService(db *gorm.DB) *TransactionService {
	return &TransactionService{DB: db}
}

func (ts *TransactionService) GetAllTransactions(w http.ResponseWriter, r *http.Request) {
	var transactions []entities.Transaction
	if err := ts.DB.Find(&transactions).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var transactionInfos []TransactionInfo
	for _, transaction := range transactions {
		transactionInfos = append(transactionInfos, TransactionInfo{
			ID:         transaction.ID,
			SenderID:   transaction.SenderID,
			ReceiverID: transaction.ReceiverID,
			Amount:     transaction.Amount,
			CreatedAt:  transaction.CreatedAt,
		})
	}

	response := SOAPEnvelope{
		Body: SOAPBody{
			GetAllTransactionsResponse: &GetAllTransactionsResponse{
				Transactions: transactionInfos,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ts *TransactionService) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the transaction request from the SOAP envelope
	req := envelope.Body.CreateTransactionRequest

	// Map the parsed data to your database model
	newTransaction := entities.Transaction{
		ID:         req.Transaction.ID,
		SenderID:   req.Transaction.SenderID,
		ReceiverID: req.Transaction.ReceiverID,
		Amount:     req.Transaction.Amount,
		CreatedAt:  time.Now(),
	}

	// Save to database
	if err := ts.DB.Create(&newTransaction).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	response := SOAPEnvelope{
		Body: SOAPBody{
			CreateTransactionResponse: &CreateTransactionResponse{
				Transaction: req.Transaction,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}
