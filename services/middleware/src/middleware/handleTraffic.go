package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"middleware/internal/db/models"

	"github.com/gofiber/fiber/v2"
)

func CheckRegisterClient(ctx *fiber.Ctx) error {
	fmt.Println("Register client")
	requestBody := ctx.Body()

	// fmt.Println(string(requestBody))

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

	fmt.Println(bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if client.Name == "" || client.Address == "" || client.BirthDate == "" || client.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	resp, err := http.Post("http://client-services:3001/api/v1/clients/register", "application/json", bytes.NewBuffer(requestBody))
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
	resp, err := http.Post("http://client-services:3001/api/v1/clients/login", "application/json", bytes.NewBuffer(requestBody))
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

	ctx.Cookie(&fiber.Cookie{
		Name:     "esb_token",
		Value:    loginResponse.Token,
		SameSite: "None",
		HTTPOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return ctx.Status(resp.StatusCode).JSON(fiber.Map{
		"message": "Response from second service",
		"token":   loginResponse.Token,
	})
}

func DoLogout(ctx *fiber.Ctx) error {
	res, err := http.Get("http://client-services:3001/api/v1/clients/logout")
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service")
	}

	defer res.Body.Close()

	// Read the response body from the second service
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service")
	}

	if res.StatusCode == http.StatusOK {
		ctx.ClearCookie("esb_token")
	}

	// Send the response from the second service back to the client
	return ctx.Status(res.StatusCode).SendString(string(body))
}

func GetAllAccounts(ctx *fiber.Ctx) error {
	//set header from cookie
	cookie := ctx.Cookies("esb_token")

	//set header from cookie
	req, err := http.NewRequest("GET", "http://account-services:3002/api/v1/accounts/", nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service: 1" + err.Error())
	}

	req.Header.Set("Cookie", "esb_token="+cookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to make request to second service: 2" + err.Error())
	}
	defer resp.Body.Close()

	// Read the response body from the second service
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service: 3" + err.Error())
	}

	// Send the response from the second service back to the client
	return ctx.Status(resp.StatusCode).SendString(string(body))
}

func CreateAccount(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("esb_token")

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
	req, err := http.NewRequest("POST", "http://account-services:3002/api/v1/accounts/", bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service")
	}

	req.Header.Set("Cookie", "esb_token="+cookie)
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

func GetAccount(ctx *fiber.Ctx) error {
	accountID := ctx.Params("id")
	cookie := ctx.Cookies("esb_token")
	// Extract the accountID from the URL
	if accountID == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("accountID is required")
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
	var account models.AccountVerify
	err := json.Unmarshal(requestBody, &account)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("Failed to unmarshal request body")
	}

	if account.Pin == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("Missing required fields")
	}

	// Make the request to the second service
	req, err := http.NewRequest("GET", "http://account-services:3002/api/v1/accounts/clientAcc/"+accountID, bytes.NewBuffer(requestBody))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service")
	}

	req.Header.Set("Cookie", "esb_token="+cookie)
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

func GetAllClientAccounts(ctx *fiber.Ctx) error {
	cookie := ctx.Cookies("esb_token")
	//set header from cookie
	req, err := http.NewRequest("GET", "http://account-services:3002/api/v1/accounts/clientAcc", nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service")
	}

	req.Header.Set("Cookie", "esb_token="+cookie)

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
	url := "http://transaction-services:3003/api/v1/transactions/getAll"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error creating request:1 " + err.Error())
	}
	req.Header.Set("Content-Type", "application/xml")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error sending request:2 " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Error reading response:3 " + err.Error())
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
	url := "http://transaction-services:3003/api/v1/transactions/getByID"
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
	url := "http://transaction-services:3003/api/v1/transactions/getByAccountID"
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
					<tns:SenderID>%s</tns:SenderID>
					<tns:ReceiverID>%s</tns:ReceiverID>
					<tns:Amount>%f</tns:Amount>
				</tns:transaction>
            </tns:CreateTransactionRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, senderID, receiverID, amount)

	// Send the XML request to the specified endpoint
	url := "http://transaction-services:3003/api/v1/transactions/create"
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
	url := "http://payment-services:3004/api/v1/payments/getAll"
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
	url := "http://payment-services:3004/api/v1/payments/getByID"
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
	url := "http://payment-services:3004/api/v1/payments/getByAccountID"
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
					<tns:AccountID>%s</tns:AccountID>
					<tns:RefCode>%s</tns:RefCode>
					<tns:Amount>%f</tns:Amount>
				</tns:payment>
            </tns:CreatePaymentRequest>
        </soapenv:Body>
    </soapenv:Envelope>`, accountID, refCode, amount)

	// Send the XML request to the specified endpoint
	url := "http://payment-services:3004/api/v1/payments/create"
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

func GetStatement(ctx *fiber.Ctx) error {
	// Extract query parameters
	accountID := ctx.Query("accountID")
	start := ctx.Query("start")
	end := ctx.Query("end")

	// Set up the URL with query parameters
	url := fmt.Sprintf("http://statement-services:3005/api/v1/statements?accountID=%s&start=%s&end=%s", accountID, start, end)

	// Create a new HTTP request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to create request to second service 1: " + err.Error())
	}

	// Manually copy headers from fasthttp.RequestHeader to http.Header
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		req.Header.Set(string(key), string(value))
	})

	// Create a new HTTP client and perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to forward request to second service 2: " + err.Error())
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("Failed to read response from second service 3: " + err.Error())
	}

	// Set the response status code and headers in the Fiber context
	ctx.Status(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			ctx.Set(key, value)
		}
	}

	// Send the response body back to the client
	return ctx.Send(body)
}
