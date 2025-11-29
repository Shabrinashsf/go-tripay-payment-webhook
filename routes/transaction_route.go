package routes

import (
	"github.com/Shabrinashsf/go-tripay-payment-webhook/controller"
	"github.com/gin-gonic/gin"
)

func Transaction(route *gin.Engine, transactionController controller.TransactionController) {
	routes := route.Group("/transaction")
	{
		routes.POST("/buy", transactionController.CreateTransaction)
		routes.POST("/webhook/tripay", transactionController.TripayWebhook)
	}
}
