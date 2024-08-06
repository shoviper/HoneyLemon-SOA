package services

import (
	"encoding/xml"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPBody struct {
	CreateTransactionRequest   *CreateTransactionRequest   `xml:"CreateTransactionRequest,omitempty"`
	CreateTransactionResponse  *CreateTransactionResponse  `xml:"CreateTransactionResponse,omitempty"`
	GetAllTransactionsResponse *GetAllTransactionsResponse `xml:"GetAllTransactionsResponse,omitempty"`

	CreatePaymentRequest   *CreatePaymentRequest   `xml:"CreatePaymentRequest,omitempty"`
	CreatePaymentResponse  *CreatePaymentResponse  `xml:"CreatePaymentResponse,omitempty"`
	GetAllPaymentsResponse *GetAllPaymentsResponse `xml:"GetAllPaymentsResponse,omitempty"`
}
