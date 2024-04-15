package main

import (
	"dogker/andrenk/billing-service/internal/rest"
	"dogker/andrenk/billing-service/internal/rest/auth"
	"dogker/andrenk/billing-service/internal/rest/deposits"
	"dogker/andrenk/billing-service/internal/rest/mutations"
	"dogker/andrenk/billing-service/internal/rest/payment"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//GORM Setup
	dsn := "host=localhost user=dogker_admin password=endeavor dbname=billing_db port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&deposits.Deposit{}, &mutations.Mutation{})

	//Public Key Setup
	publicKeyPEM := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEnlwXdOFOQFhhEoYksncm/mmRMjVv\nVKiJhzabtB5d2uMV7Xn0SKVzJB4jKUM/05Qcfmxkjt4OyBJNQ4LE5oa3eQ==\n-----END PUBLIC KEY-----\n"

	err = auth.InitPublicKey(publicKeyPEM)
	if err != nil {
		fmt.Println("Failed to initialize public key:", err)
		return
	}

	//Auth Setup
	authService := auth.NewService()

	//Payment Setup
	paymentService := payment.NewService()

	//Mutation Setup
	mutationRepository := mutations.NewRepository(db)
	mutationService := mutations.NewService(mutationRepository)

	//Deposits Setup
	depositRepository := deposits.NewRepository(db)
	depositService := deposits.NewService(depositRepository)
	depositsHandler := deposits.NewDepositsHandler(depositService, paymentService, authService, mutationService)

	//REST Setup
	router := gin.Default()
	api := router.Group("/api/v1")

	depositsRoute := api.Group("/deposits")
	{
		depositsRoute.GET("/", rest.AuthMiddleware(), depositsHandler.GetDepositByID)
		depositsRoute.POST("/", rest.AuthMiddleware(), depositsHandler.InitiateDeposit)
	}

	notificationRoute := api.Group("/notification")
	{
		notificationRoute.POST("/", depositsHandler.HandleNotificationPayment)
	}

	router.Run(":8080")
}
