package dto

import "errors"

const (
	MESSAGE_SUCCESS_CREATE_TRANSACTION  = "success create transaction"
	MESSAGE_SUCCESS_GET_CALLBACK_TRIPAY = "success get callback tripay"

	MESSAGE_FAILED_GET_DATA_FROM_BODY  = "failed to get data from body"
	MESSAGE_FAILED_CREATE_TRANSACTION  = "failed to create transaction"
	MESSAGE_FAILED_GET_CALLBACK_TRIPAY = "failed get callback tripay"
)

var (
	ErrUnrecognizedCallbackEvent  = errors.New("unrecognized callback event")
	ErrInvalidSignature           = errors.New("invalid signature")
	ErrOnlyClosedPaymentSupported = errors.New("only closed payment supported")
	ErrTransactionNotFound        = errors.New("transaction not found")
	ErrUnknownStatus              = errors.New("unknown status")
	ErrFailedToUpdateStatus       = errors.New("failed to update transaction status")
)

type (
	CreateTransactionRequest struct {
		Name          string `json:"name" binding:"required"`
		Email         string `json:"email" binding:"required,email"`
		MobileNumber  string `json:"mobile_number" binding:"required"`
		ProductID     string `json:"product_id" binding:"required,uuid"`
		PaymentMethod string `json:"payment_method" binding:"required"`
	}

	CreateTransactionResponse struct {
		InvoiceURL string `json:"invoice_url"`
	}

	TripayExpiredTime int

	TripayOrderRequest struct {
		Method        string                    `json:"method"`
		MerchantRef   string                    `json:"merchant_ref"`
		Amount        int                       `json:"amount"`
		CustomerName  string                    `json:"customer_name"`
		CustomerEmail string                    `json:"customer_email"`
		CustomerPhone string                    `json:"customer_phone"`
		OrderItems    []OrderItemPaymentRequest `json:"order_items"`
		ReturnURL     string                    `json:"return_url"`
		ExpiredTime   TripayExpiredTime         `json:"expired_time"`
		Signature     string                    `json:"signature"`
	}

	OrderItemPaymentRequest struct {
		SKU        string `json:"sku"`
		Name       string `json:"name"`
		Price      int    `json:"price"`
		Quantity   int    `json:"quantity"`
		ProductURL string `json:"product_url"`
		ImageURL   string `json:"image_url"`
	}

	TripayResponse struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
		Data    Data   `json:"data"`
	}

	Data struct {
		Reference   string `json:"reference"`
		MerchantRef string `json:"merchant_ref"`
		PaymentURL  string `json:"checkout_url"`
	}

	TripayWebhookRequest struct {
		Reference         string `json:"reference"`
		MerchantRef       string `json:"merchant_ref"`
		PaymentMethod     string `json:"payment_method"`
		PaymentMethodCode string `json:"payment_method_code"`
		TotalAmount       int    `json:"total_amount"`
		FeeMerchant       int    `json:"fee_merchant"`
		FeeCustomer       int    `json:"fee_customer"`
		TotalFee          int    `json:"total_fee"`
		AmountReceived    int    `json:"amount_received"`
		IsClosedPayment   int    `json:"is_closed_payment"`
		Status            string `json:"status"`
		PaidAt            int    `json:"paid_at"`
	}

	TripayWebhookResponse struct {
		Success bool `json:"success"`
	}
)
