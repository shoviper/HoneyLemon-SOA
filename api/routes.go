package api

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"soaProject/api/services"

	"github.com/gofiber/fiber/v2"
	viper "github.com/spf13/viper"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
)

func convertURI(uri *fasthttp.URI) *url.URL {
	return &url.URL{
		Scheme:   string(uri.Scheme()),
		Host:     string(uri.Host()),
		Path:     string(uri.Path()),
		RawQuery: string(uri.QueryString()),
	}
}

func bytesToReadCloser(b []byte) io.ReadCloser {
	return io.NopCloser(bytes.NewReader(b))
}

func WrapHandler(handler http.HandlerFunc) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		req := c.Request()
		r := &http.Request{
			Method: string(req.Header.Method()),
			URL:    convertURI(req.URI()),
			Proto:  "HTTP/1.1",
			Header: make(http.Header),
			Body:   bytesToReadCloser(req.Body()),
		}

		// Copy headers
		req.Header.VisitAll(func(k, v []byte) {
			r.Header.Add(string(k), string(v))
		})

		writer := &fiberResponseWriter{c.Response()}

		handler(writer, r)

		return nil
	}
}

type fiberResponseWriter struct {
	resp *fasthttp.Response
}

func (w *fiberResponseWriter) Header() http.Header {
	header := http.Header{}
	w.resp.Header.VisitAll(func(key, value []byte) {
		header.Set(string(key), string(value))
	})
	return header
}

func (w *fiberResponseWriter) Write(b []byte) (int, error) {
	w.resp.SetBody(b)
	return len(b), nil
}

func (w *fiberResponseWriter) WriteHeader(statusCode int) {
	w.resp.SetStatusCode(statusCode)
}

func SetupRoutes(app *fiber.App, db *gorm.DB, vp *viper.Viper) {
	// clientService := services.NewClientService(db, vp)
	// accountService := services.NewAccountService(db, vp)
	transactionService := services.NewTransactionService(db)
	paymentService := services.NewPaymentService(db)

	// JWTServ := services.NewJWTConfig(vp)

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	api := app.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			// client := v1.Group("/clients")
			// {
			// 	client.Get("/", clientService.GetAllClients)
			// 	client.Post("/register", clientService.RegisterClient)
			// 	client.Post("/login", clientService.LoginClient)
			// 	client.Delete("/", clientService.DeleteClient)

			// }
			// account := v1.Group("/accounts")
			// {
			// 	account.Use(JWTServ.JWTAuth())
			// 	account.Get("/", accountService.GetAllAccounts)
			// 	account.Post("/", accountService.CreateAccount)
			// 	account.Get("/clientAcc/:id", accountService.GetAccount)
			// 	account.Get("/clientAcc", accountService.GetAllClientAccounts)
			// }
			transaction := v1.Group("/transactions")
			{
				transaction.Post("/getAll", services.WrapHandler(transactionService.GetAllTransactions))
				transaction.Post("/getByID", services.WrapHandler(transactionService.GetTransactionByID))
				transaction.Post("/getByAccountID", services.WrapHandler(transactionService.GetTransactionsByAccountID))
				transaction.Post("/create", services.WrapHandler(transactionService.CreateTransaction))
			}
			payment := v1.Group("/payments")
			{
				payment.Post("/getAll", services.WrapHandler(paymentService.GetAllPayments))
				payment.Post("/getByID", services.WrapHandler(paymentService.GetPaymentByID))
				payment.Post("/getByAccountID", services.WrapHandler(paymentService.GetPaymentsByAccountID))
				payment.Post("/create", services.WrapHandler(paymentService.CreatePayment))
			}
		}
	}
}
