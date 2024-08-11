package payment

import (
	"encoding/xml"
	"net/http"
	"payment/internal/db/entities"
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

type GetPaymentByIDRequest struct {
	XMLName   xml.Name `xml:"GetPaymentByIDRequest"`
	PaymentID string   `xml:"PaymentID"`
}

type GetPaymentByIDResponse struct {
	Payments PaymentInfo `xml:"payment"`
}

type GetPaymentsByAccountIDRequest struct {
	XMLName   xml.Name `xml:"GetPaymentsByAccountIDRequest"`
	AccountID string   `xml:"AccountID"`
}

type GetPaymentsByAccountIDResponse struct {
	Payments []PaymentInfo `xml:"payment"`
}

type NullPaymentResponse struct {
	Payments []string `xml:"payments"`
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

func (ps *PaymentService) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	var payments []entities.Payment
	if err := ps.DB.Find(&payments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(payments) == 0 {
		response := SOAPEnvelope{
			Body: SOAPBody{
				NullPaymentResponse: &NullPaymentResponse{
					Payments: []string{""},
				},
			},
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
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
			GetAllPaymentsResponse: &GetAllPaymentsResponse{
				Payments: paymentInfos,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ps *PaymentService) GetPaymentByID(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the payment request from the SOAP envelope
	req := envelope.Body.GetPaymentByIDRequest
	if req == nil {
		http.Error(w, "Invalid request: missing GetPaymentByIDRequest", http.StatusBadRequest)
		return
	}

	var payment entities.Payment
	if err := ps.DB.Where("id = ?", req.PaymentID).Find(&payment).Error; err != nil || payment.ID == "" {
		response := SOAPEnvelope{
			Body: SOAPBody{
				NullPaymentResponse: &NullPaymentResponse{
					Payments: []string{""},
				},
			},
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
		return
	}

	var paymentInfo = PaymentInfo{
		ID:        payment.ID,
		AccountID: payment.AccountID,
		RefCode:   payment.RefCode,
		Amount:    payment.Amount,
		CreatedAt: payment.CreatedAt,
	}

	response := SOAPEnvelope{
		Body: SOAPBody{
			GetPaymentByIDResponse: &GetPaymentByIDResponse{
				Payments: paymentInfo,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ps *PaymentService) GetPaymentsByAccountID(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the payment request from the SOAP envelope
	req := envelope.Body.GetPaymentsByAccountIDRequest

	var payments []entities.Payment
	if err := ps.DB.Where("account_id = ?", req.AccountID).Find(&payments).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if len(payments) == 0 {
		response := SOAPEnvelope{
			Body: SOAPBody{
				NullPaymentResponse: &NullPaymentResponse{
					Payments: []string{""},
				},
			},
		}
		w.Header().Set("Content-Type", "application/xml")
		xml.NewEncoder(w).Encode(response)
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
			GetPaymentsByAccountIDResponse: &GetPaymentsByAccountIDResponse{
				Payments: paymentInfos,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}

func (ps *PaymentService) CreatePayment(w http.ResponseWriter, r *http.Request) {
	var envelope SOAPEnvelope
	if err := xml.NewDecoder(r.Body).Decode(&envelope); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Extract the payment request from the SOAP envelope
	req := envelope.Body.CreatePaymentRequest

	// Start a database payment
	tx := ps.DB.Begin()

	// Retrieve account and receiver accounts
	var account entities.Account
	if err := tx.Where("id = ?", req.Payment.AccountID).First(&account).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Account not found", http.StatusBadRequest)
		return
	}

	// Check if the account has sufficient balance
	if account.Balance < req.Payment.Amount {
		tx.Rollback()
		http.Error(w, "Insufficient balance", http.StatusBadRequest)
		return
	}

	// Update the account's and receiver's balances
	account.Balance -= req.Payment.Amount

	// Save the updated accounts back to the database
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update account's balance", http.StatusInternalServerError)
		return
	}

	// Map the parsed data to your database model
	newPayment := entities.Payment{
		AccountID: req.Payment.AccountID,
		RefCode:   req.Payment.RefCode,
		Amount:    req.Payment.Amount,
	}

	// Save to database
	if err := ps.DB.Create(&newPayment).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Use the saved payment to include ID and timestamps in the response
	req.Payment.ID = newPayment.ID
	req.Payment.CreatedAt = newPayment.CreatedAt

	// Commit the Payment
	tx.Commit()

	// Respond with success
	response := SOAPEnvelope{
		Body: SOAPBody{
			CreatePaymentResponse: &CreatePaymentResponse{
				Payment: req.Payment,
			},
		},
	}

	w.Header().Set("Content-Type", "application/xml")
	xml.NewEncoder(w).Encode(response)
}
