package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/dto"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/service"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/utils/response"
	"github.com/gin-gonic/gin"
)

type (
	TransactionController interface {
		CreateTransaction(ctx *gin.Context)
		TripayWebhook(ctx *gin.Context)
	}

	transactionController struct {
		transactionService service.TransactionService
	}
)

func NewTransactionController(transactionService service.TransactionService) TransactionController {
	return &transactionController{
		transactionService: transactionService,
	}
}

func (c *transactionController) CreateTransaction(ctx *gin.Context) {
	// Dont forget to implement your business logic
	var req dto.CreateTransactionRequest
	if err := ctx.ShouldBind(&req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.JSON(http.StatusBadRequest, res)
		return
	}

	result, err := c.transactionService.CreateTransaction(ctx.Request.Context(), req)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_CREATE_TRANSACTION, err.Error(), nil)
		ctx.JSON(http.StatusInternalServerError, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_CREATE_TRANSACTION, result)
	ctx.JSON(http.StatusOK, res)
}

func (c *transactionController) TripayWebhook(ctx *gin.Context) {
	svcCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// 1. Ambil raw body sekali
	rawBody, err := ctx.GetRawData()
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// 2. Parse JSON ke struct
	var req dto.TripayWebhookRequest
	if err := json.Unmarshal(rawBody, &req); err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_DATA_FROM_BODY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	// 3. Ambil header
	cbSignature := ctx.GetHeader("X-Callback-Signature")
	cbEvent := ctx.GetHeader("X-Callback-Event")

	// 4. Kirim semuanya ke service
	_, err = c.transactionService.TripayWebhook(svcCtx, rawBody, req, cbSignature, cbEvent)
	if err != nil {
		res := response.BuildResponseFailed(dto.MESSAGE_FAILED_GET_CALLBACK_TRIPAY, err.Error(), nil)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, res)
		return
	}

	res := response.BuildResponseSuccess(dto.MESSAGE_SUCCESS_GET_CALLBACK_TRIPAY, nil)
	ctx.JSON(http.StatusOK, res)
}
