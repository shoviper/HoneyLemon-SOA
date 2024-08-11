package payment

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type SOAPEnvelope struct {
	XMLName xml.Name `xml:"http://schemas.xmlsoap.org/soap/envelope/ Envelope"`
	Body    SOAPBody `xml:"Body"`
}

type SOAPHeader struct{}

type SOAPBody struct {
	CreatePaymentRequest           *CreatePaymentRequest           `xml:"CreatePaymentRequest,omitempty"`
	CreatePaymentResponse          *CreatePaymentResponse          `xml:"CreatePaymentResponse,omitempty"`
	GetAllPaymentsResponse         *GetAllPaymentsResponse         `xml:"GetAllPaymentsResponse,omitempty"`
	GetPaymentByIDRequest          *GetPaymentByIDRequest          `xml:"GetPaymentByIDRequest,omitempty"`
	GetPaymentByIDResponse         *GetPaymentByIDResponse         `xml:"GetPaymentByIDResponse,omitempty"`
	GetPaymentsByAccountIDRequest  *GetPaymentsByAccountIDRequest  `xml:"GetPaymentsByAccountIDRequest,omitempty"`
	GetPaymentsByAccountIDResponse *GetPaymentsByAccountIDResponse `xml:"GetPaymentsByAccountIDResponse,omitempty"`
	NullPaymentResponse            *NullPaymentResponse            `xml:"NullPaymentResponse,omitempty"`
}

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
