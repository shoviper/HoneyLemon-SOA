package services

import (
	"encoding/xml"
	"net/http"
	"soaProject/internal/db/entities"
	"time"

	"gorm.io/gorm"
)

type CreatePaymentRequest struct {
	XMLName xml.Name    `xml:"CreatePaymentRequest"`
	Payment PaymentInfo `xml:"payment"`
}

type CreatePaymentResponse struct {
	Payment PaymentInfo `xml:"payment"`
}

type GetAllPaymentsResponse struct {
	Payments []PaymentInfo `xml:"payments"`
}

type PaymentInfo struct {
	ID        string    `xml:"ID"`
	AccountID string    `xml:"AccountID"`
	RefCode   string    `xml:"RefCode"`
	Amount    float64   `xml:"Amount"`
	CreatedAt time.Time `xml:"CreatedAt"`
}

type PaymentService struct {
	DB *gorm.DB
}

func NewPaymentService(db *gorm.DB) *PaymentService {
	return &PaymentService{DB: db}
}

func (ts *PaymentService) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	var payments []entities.Payment
	if err := ts.DB.Find(&payments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var paymentInfos []PaymentInfo
	for _, payment := range payments {
		paymentInfos = append(paymentInfos, PaymentInfo{
			ID:        payment.ID,
			AccountID: payment.AccountID,
			RefCode:   payment.RefCode,
			Amount:    payment.Amount,
			CreatedAt: payment.CreatedAt,
		})
	}

	response := SOAPEnvelope{
		Body: SOAPBody{
			GetAllPaymentsResponse: GetAllPaymentsResponse{
				Payments: paymentInfos,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ts *PaymentService) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the payment request from the SOAP envelope
	req := envelope.Body.CreatePaymentRequest

	// Map the parsed data to your database model
	newTransaction := entities.Payment{
		ID:        req.Payment.ID,
		AccountID: req.Payment.AccountID,
		RefCode:   req.Payment.RefCode,
		Amount:    req.Payment.Amount,
		CreatedAt: time.Now(),
	}

	// Save to database
	if err := ts.DB.Create(&newTransaction).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Respond with success
	response := SOAPEnvelope{
		Body: SOAPBody{
			CreatePaymentResponse: CreatePaymentResponse{
				Payment: req.Payment,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}
