package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Shabrinashsf/go-tripay-payment-webhook/cmd"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/config"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/controller"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/middleware"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/repository"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/routes"
	"github.com/Shabrinashsf/go-tripay-payment-webhook/service"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Println("cannot load .env:", err)
	} else {
		log.Println(".env loaded successfully")
	}
}

func main() {
	db := config.SetUpDatabaseConnection()
	defer config.CloseDatabaseConnection(db)

	if len(os.Args) > 1 {
		cmd.Commands(db)
		return
	}

	var (
		// Implementation Dependency Injection
		// Repository
		transactionRepository repository.TransactionRepository = repository.NewTransactionRepository(db)

		// Service
		transactionService service.TransactionService = service.NewTransactionService(transactionRepository)

		// Controller
		transactionController controller.TransactionController = controller.NewTransactionController(transactionService)
	)

	server := gin.Default()
	server.Use(middleware.CORSMiddleware())

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  http.StatusNotFound,
			"message": "Route Not Found",
		})
	})

	server.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// routes
	routes.Transaction(server, transactionController)

	server.Static("/assets", "./assets")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}

	var serve string
	if os.Getenv("APP_ENV") == "localhost" {
		serve = "127.0.0.1:" + port
	} else {
		serve = ":" + port
	}

	if err := server.Run(serve); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}
