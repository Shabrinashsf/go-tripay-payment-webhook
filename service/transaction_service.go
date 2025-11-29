package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/dto"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/entity"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/repository"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/utils/payment"
	"github.com/google/uuid"
)

type (
	TransactionService interface {
		CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (dto.CreateTransactionResponse, error)
		TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error)
	}

	transactionService struct {
		transactionRepo repository.TransactionRepository
	}
)

func NewTransactionService(transactionRepo repository.TransactionRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
	}
}

func (s *transactionService) CreateTransaction(ctx context.Context, req dto.CreateTransactionRequest) (dto.CreateTransactionResponse, error) {
	// Dont forget to implement your business logic
	product, err := s.transactionRepo.GetProductByID(ctx, nil, uuid.MustParse(req.ProductID))
	if err != nil {
		return dto.CreateTransactionResponse{}, err
	}

	// build Order Request
	merchantRef := "INV-" + product.ID.String()

	orderItems := []dto.OrderItemPaymentRequest{
		{
			SKU:      product.Name,
			Name:     product.Name,
			Price:    int(product.Price),
			Quantity: 1,
		},
	}

	invoice := dto.TripayOrderRequest{
		Method:        req.PaymentMethod,
		MerchantRef:   merchantRef,
		Amount:        int(product.Price),
		CustomerName:  req.Name,
		CustomerEmail: req.Email,
		CustomerPhone: req.MobileNumber,
		OrderItems:    orderItems,
		ReturnURL:     os.Getenv("RETURN_URL"),
		ExpiredTime:   dto.TripayExpiredTime(time.Now().Unix() + (60 * 60 * 24)), // 24 jam
	}

	// bikin invoice
	url, err := payment.SendTripayInvoice(ctx, invoice)

	if err != nil {
		panic(err)
	}

	fmt.Println("Checkout URL:", url)

	// bikin transaksi
	trans := entity.Transaction{
		ProductID:  product.ID,
		AmountPaid: 0,
		Status:     "PENDING",
		InvoiceUrl: url.Data.PaymentURL,
		Reference:  url.Data.Reference,
	}

	// setor ke db
	transaction, err := s.transactionRepo.CreateTransaction(ctx, nil, trans)
	if err != nil {
		return dto.CreateTransactionResponse{}, err
	}

	return dto.CreateTransactionResponse{
		InvoiceURL: transaction.InvoiceUrl,
	}, nil
}

func (s *transactionService) TripayWebhook(ctx context.Context, rawBody []byte, payload dto.TripayWebhookRequest, callbackSignature string, event string) (dto.TripayWebhookResponse, error) {
	privateKey := os.Getenv("TRIPAY_PRIVATE_KEY")

	if event != "payment_status" {
		return dto.TripayWebhookResponse{}, dto.ErrUnrecognizedCallbackEvent
	}

	// hitung signature local
	mac := hmac.New(sha256.New, []byte(privateKey))
	mac.Write(rawBody)
	localSignature := hex.EncodeToString(mac.Sum(nil))

	if localSignature != callbackSignature {
		return dto.TripayWebhookResponse{}, dto.ErrInvalidSignature
	}

	// pastikan closed payment
	if payload.IsClosedPayment != 1 {
		return dto.TripayWebhookResponse{}, dto.ErrOnlyClosedPaymentSupported
	}

	// cari transaksi berdasarkan reference
	transaction, err := s.transactionRepo.GetTransactionByReference(ctx, nil, payload.Reference)
	if err != nil {
		return dto.TripayWebhookResponse{}, dto.ErrTransactionNotFound
	}

	// update status transaksi
	switch strings.ToUpper(payload.Status) {
	case "PAID":
		transaction.Status = "PAID"
	case "FAILED":
		transaction.Status = "FAILED"
	case "EXPIRED":
		transaction.Status = "EXPIRED"
	case "REFUND":
		transaction.Status = "REFUND"
	default:
		return dto.TripayWebhookResponse{}, dto.ErrUnknownStatus
	}

	transaction.AmountPaid = payload.TotalAmount

	// simpan ke database
	if err := s.transactionRepo.UpdateTransaction(ctx, nil, transaction); err != nil {
		return dto.TripayWebhookResponse{}, dto.ErrFailedToUpdateStatus
	}

	return dto.TripayWebhookResponse{
		Success: true,
	}, nil
}
