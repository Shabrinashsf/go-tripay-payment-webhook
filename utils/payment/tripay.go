package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/dto"
)

type (
	Signature struct {
		Amount       int64
		PrivateKey   string
		MerchantCode string
		MerchanReff  string
		Channel      string
	}

	Client struct {
		MerchantCode string
		ApiKey       string
		PrivateKey   string
		Mode         string
		signature    Signature
	}
)

func (s *Signature) CreateSignature() string {
	var signStr string
	if s.Amount != 0 {
		signStr = s.MerchantCode + s.MerchanReff + fmt.Sprintf("%d", s.Amount)
	} else {
		signStr = s.MerchantCode + s.Channel + s.MerchanReff
	}

	key := []byte(s.PrivateKey)
	message := []byte(signStr)

	hash := hmac.New(sha256.New, key)
	hash.Write(message)
	signature := hex.EncodeToString(hash.Sum(nil))

	return signature
}

func (c *Client) SetSignature(sig Signature) {
	c.signature = sig
}

func (c Client) BaseUrl() string {
	if c.Mode == "development" {
		return "https://tripay.co.id/api-sandbox"
	}

	return "https://tripay.co.id/api"
}

func (c *Client) CreateTransaction(ctx context.Context, req dto.TripayOrderRequest) (dto.TripayResponse, error) {
	if c.signature.MerchanReff == "" {
		return dto.TripayResponse{}, errors.New("signature not set")
	}

	// Tambah signature ke request
	requestBody := map[string]interface{}{
		"method":         req.Method,
		"merchant_ref":   req.MerchantRef,
		"amount":         req.Amount,
		"customer_name":  req.CustomerName,
		"customer_email": req.CustomerEmail,
		"customer_phone": req.CustomerPhone,
		"order_items":    req.OrderItems,
		"expired_time":   int64(req.ExpiredTime),
		"return_url":     req.ReturnURL,
		"signature":      c.signature.CreateSignature(),
	}

	jsonBody, _ := json.Marshal(requestBody)

	url := c.BaseUrl() + "/transaction/create"

	httpReq, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonBody))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.ApiKey)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return dto.TripayResponse{}, err
	}

	defer resp.Body.Close()
	bodyBytes, _ := io.ReadAll(resp.Body)

	var parsed dto.TripayResponse
	if err := json.Unmarshal(bodyBytes, &parsed); err != nil {
		return dto.TripayResponse{}, err
	}

	if !parsed.Success {
		return dto.TripayResponse{}, errors.New(parsed.Message)
	}

	return parsed, nil
}

func SendTripayInvoice(ctx context.Context, invoice dto.TripayOrderRequest) (dto.TripayResponse, error) {
	// 1. Create Signature
	sig := Signature{
		Amount:       int64(invoice.Amount),
		PrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		MerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		MerchanReff:  invoice.MerchantRef,
	}

	// 2. Create Client
	client := Client{
		MerchantCode: os.Getenv("TRIPAY_MERCHANT_CODE"),
		ApiKey:       os.Getenv("TRIPAY_API_KEY"),
		PrivateKey:   os.Getenv("TRIPAY_PRIVATE_KEY"),
		Mode:         "development",
	}

	// 3. Set Signature
	client.SetSignature(sig)

	// 4. Call Tripay API
	res, err := client.CreateTransaction(ctx, invoice)
	if err != nil {
		return dto.TripayResponse{}, err
	}

	// Return only checkout URL
	return res, nil
}
