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

type GetTransactionByIDRequest struct {
	XMLName       xml.Name `xml:"GetTransactionByIDRequest"`
	TransactionID string   `xml:"TransactionID"`
}

type GetTransactionByIDResponse struct {
	Transactions TransactionInfo `xml:"transactions"`
}

type GetTransactionsByAccountIDRequest struct {
	XMLName   xml.Name `xml:"GetTransactionsByAccountIDRequest"`
	AccountID string   `xml:"AccountID"`
}

type GetTransactionsByAccountIDResponse struct {
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

func (ts *TransactionService) GetTransactionByID(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the transaction request from the SOAP envelope
	req := envelope.Body.GetTransactionByIDRequest
	if req == nil {
		http.Error(w, "Invalid request: missing GetTransactionByIDRequest", http.StatusBadRequest)
		return
	}

	var transactions []entities.Transaction
	if err := ts.DB.Where("id = ?", req.TransactionID).Find(&transactions).Error; err != nil {
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
			GetTransactionsByAccountIDResponse: &GetTransactionsByAccountIDResponse{
				Transactions: transactionInfos,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ts *TransactionService) GetTransactionsByAccountID(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the transaction request from the SOAP envelope
	req := envelope.Body.GetTransactionsByAccountIDRequest

	var transactions []entities.Transaction
	if err := ts.DB.Where("sender_id = ? OR receiver_id = ?", req.AccountID, req.AccountID).Find(&transactions).Error; err != nil {
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
			GetTransactionsByAccountIDResponse: &GetTransactionsByAccountIDResponse{
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

	// Start a database transaction
	tx := ts.DB.Begin()

	// Retrieve sender and receiver accounts
	var sender entities.Account
	if err := tx.Where("id = ?", req.Transaction.SenderID).First(&sender).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Sender account not found", http.StatusBadRequest)
		return
	}

	var receiver entities.Account
	if err := tx.Where("id = ?", req.Transaction.ReceiverID).First(&receiver).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Receiver account not found", http.StatusBadRequest)
		return
	}

	// Check if the sender has sufficient balance
	if sender.Balance < req.Transaction.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Update the sender's and receiver's balances
	sender.Balance -= req.Transaction.Amount
	receiver.Balance += req.Transaction.Amount

	// Save the updated accounts back to the database
	if err := tx.Save(&sender).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update sender's balance", http.StatusInternalServerError)
		return
	}

	if err := tx.Save(&receiver).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update receiver's balance", http.StatusInternalServerError)
		return
	}

	// Create and save the new transaction
	newTransaction := entities.Transaction{
		ID:         req.Transaction.ID,
		SenderID:   req.Transaction.SenderID,
		ReceiverID: req.Transaction.ReceiverID,
		Amount:     req.Transaction.Amount,
		CreatedAt:  time.Now(),
	}

	if err := tx.Create(&newTransaction).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create transaction", http.StatusInternalServerError)
		return
	}

	// Commit the transaction
	tx.Commit()

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
