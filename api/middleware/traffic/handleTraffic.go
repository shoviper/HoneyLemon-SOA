package traffic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"soaProject/internal/db/models"
	local "soaProject/internal/local"

	"github.com/gofiber/fiber/v2"
)

func CheckRegisterClient(ctx *fiber.Ctx) error {
	requestBody := ctx.Body()

	//check if the request body is empty
	if len(requestBody) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is empty")
	}

	//check if the request body is valid JSON
	if !json.Valid(requestBody) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is not valid JSON")
	}

	//check input fields
	var client models.RegisterClient
	err := json.Unmarshal(requestBody, &client)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if client.Name == "" || client.Address == "" || client.BirthDate == "" || client.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	resp, err := http.Post("http://localhost:3000/api/v1/clients/register", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()

	// Read the response body from the second service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).JSON(fiber.Map{
		"message": "Response from second service",
		"body":    string(body),
	})
}

func CheckLoginClient(ctx *fiber.Ctx) error {
	requestBody := ctx.Body()

	//check if the request body is empty
	if len(requestBody) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is empty")
	}

	//check if the request body is valid JSON
	if !json.Valid(requestBody) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is not valid JSON")
	}

	//check input fields
	var client models.LoginClient
	err := json.Unmarshal(requestBody, &client)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if client.ID == "" || client.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	resp, err := http.Post("http://localhost:3000/api/v1/clients/login", "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()

	// Read the response body from the second service in JSON
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	var loginResponse models.LoginResponse
	if err := json.Unmarshal(body, &loginResponse); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to parse JSON response from second service")
	}

	randomStr := local.GenerateRandomString(10)
	// Set a linked token in the ESB service
	linkedToken := randomStr+"."+ loginResponse.Token

	ctx.Cookie(&fiber.Cookie{
		Name:     "esb_token",
		Value:    linkedToken,
		SameSite: "None",
		Expires: time.Now().Add(24 * time.Hour),
	})
	
	return ctx.Status(resp.StatusCode).JSON(fiber.Map{
		"message": "Response from second service",
		"token":   linkedToken,
	})
}

func GetAllAccounts(ctx *fiber.Ctx) error {
	//set header from cookie
	cookie := ctx.Cookies("backend_token")
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Empty backend token")
	}

	//set header from cookie
	req, err := http.NewRequest("GET", "http://localhost:3000/api/v1/accounts/", nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service")
	}

	req.Header.Set("Cookie", "backend_token="+cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()


	// Read the response body from the second service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).SendString(string(body))
}

func CreateAccount(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("backend_token")
	if cookie == "" {
		return ctx.Status(fiber.StatusUnauthorized).SendString("Empty backend token")
	}

	requestBody := ctx.Body()

	//check if the request body is empty
	if len(requestBody) == 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is empty")
	}

	//check if the request body is valid JSON
	if !json.Valid(requestBody) {
		return ctx.Status(fiber.StatusBadRequest).SendString("Request body is not valid JSON")
	}

	//check input fields
	var account models.CreateAccount
	err := json.Unmarshal(requestBody, &account)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if account.Type == "" || account.Pin == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	req, err := http.NewRequest("POST", "http://localhost:3000/api/v1/accounts/", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service")
	}

	req.Header.Set("Cookie", "backend_token="+cookie)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}
	defer resp.Body.Close()

	// Read the response body from the second service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).SendString(string(body))
}

func GetAllTransactions(ctx *fiber.Ctx) error {
	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
   		<soapenv:Header/>
   		<soapenv:Body>
      		<tns:GetAllTransactionsRequest/>
   		</soapenv:Body>
	</soapenv:Envelope>`)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/transactions/getAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetAllTransactionsResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func GetTransactionByID(ctx *fiber.Ctx) error {
	// Extract query parameter
	transactionID := ctx.Query("transactionID")
	if transactionID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("transactionID is required")
	}

	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:GetTransactionByIDRequest>
                <tns:TransactionID>%s</tns:TransactionID>
            </tns:GetTransactionByIDRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, transactionID)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/transactions/getByID"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetTransactionByIDResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func GetTransactionsByAccountID(ctx *fiber.Ctx) error {
	// Extract query parameter
	accountID := ctx.Query("accountID")
	if accountID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("accountID is required")
	}

	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:GetTransactionsByAccountIDRequest>
                <tns:AccountID>%s</tns:AccountID>
            </tns:GetTransactionsByAccountIDRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, accountID)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/transactions/getByAccountID"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetTransactionsByAccountIDResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func CreateTransaction(ctx *fiber.Ctx) error {
	// Read the raw JSON body
	rawBody := ctx.Body()

	// Parse the JSON body
	var requestBody map[string]interface{}
	if err := json.Unmarshal(rawBody, &requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Error parsing request body: " + err.Error())
	}

	// Extract values from the parsed JSON
	transactionID, ok := requestBody["transactionID"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid transactionID")
	}
	senderID, ok := requestBody["senderID"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid senderID")
	}
	receiverID, ok := requestBody["receiverID"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid receiverID")
	}
	amount, ok := requestBody["amount"].(float64) // JSON numbers are parsed as float64
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid amount")
	}

	// Create the request body for the SOAP request
	requestBodyXML := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:CreateTransactionRequest>
				<tns:transaction>
					<tns:ID>%s</tns:ID>
					<tns:SenderID>%s</tns:SenderID>
					<tns:ReceiverID>%s</tns:ReceiverID>
					<tns:Amount>%f</tns:Amount>
				</tns:transaction>
            </tns:CreateTransactionRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, transactionID, senderID, receiverID, amount)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/transactions/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBodyXML)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "CreateTransactionResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func GetAllPayments(ctx *fiber.Ctx) error {
	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
		<soapenv:Header/>
		<soapenv:Body>
			<tns:GetAllPaymentRequest/>
		</soapenv:Body>
	</soapenv:Envelope>`)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/payments/getAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetAllPaymentsResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func GetPaymentByID(ctx *fiber.Ctx) error {
	// Extract query parameter
	paymentID := ctx.Query("paymentID")
	if paymentID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("paymentID is required")
	}

	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:GetPaymentByIDRequest>
                <tns:PaymentID>%s</tns:PaymentID>
            </tns:GetPaymentByIDRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, paymentID)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/payments/getByID"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetPaymentByIDResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func GetPaymentsByAccountID(ctx *fiber.Ctx) error {
	// Extract query parameter
	accountID := ctx.Query("accountID")
	if accountID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("accountID is required")
	}

	// Create the request body for the SOAP request
	requestBody := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:GetPaymentsByAccountIDRequest>
                <tns:AccountID>%s</tns:AccountID>
            </tns:GetPaymentsByAccountIDRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, accountID)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/payments/getByAccountID"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "GetPaymentsByAccountIDResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}

func CreatePayment(ctx *fiber.Ctx) error {
	// Read the raw JSON body
	rawBody := ctx.Body()

	// Parse the JSON body
	var requestBody map[string]interface{}
	if err := json.Unmarshal(rawBody, &requestBody); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Error parsing request body: " + err.Error())
	}

	// Extract values from the parsed JSON
	paymentID, ok := requestBody["paymentID"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid transactionID")
	}
	accountID, ok := requestBody["accountID"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid senderID")
	}
	refCode, ok := requestBody["refCode"].(string)
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid receiverID")
	}
	amount, ok := requestBody["amount"].(float64) // JSON numbers are parsed as float64
	if !ok {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing or invalid amount")
	}

	// Create the request body for the SOAP request
	requestBodyXML := fmt.Sprintf(`
    <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:tns="http://soaProject/TransactionService">
        <soapenv:Header/>
        <soapenv:Body>
            <tns:CreatePaymentRequest>
				<tns:payment>
					<tns:ID>%s</tns:ID>
					<tns:AccountID>%s</tns:AccountID>
					<tns:RefCode>%s</tns:RefCode>
					<tns:Amount>%f</tns:Amount>
				</tns:payment>
            </tns:CreatePaymentRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, paymentID, accountID, refCode, amount)

	// Send the XML request to the specified endpoint
	url := "http://localhost:3000/api/v1/payments/create"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBodyXML)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request: " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response: " + err.Error())
	}

	// Convert XML response to JSON
	jsonResponse, err := ConvertXMLToJSON(responseBody)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error converting XML to JSON: " + err.Error())
	}

	// Extract the Body part and replace "content"
	bodyJson, err := ExtractBody(jsonResponse, "CreatePaymentResponse")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error extracting Body from JSON: " + err.Error())
	}

	// Return the response back to the client
	return ctx.Status(resp.StatusCode).SendString(string(bodyJson))
}
