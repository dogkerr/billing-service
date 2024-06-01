package main

import (
	"dogker/andrenk/billing-service/internal/rabbitmq"
	"dogker/andrenk/billing-service/internal/rest"
	"dogker/andrenk/billing-service/internal/rest/auth"
	"dogker/andrenk/billing-service/internal/rest/charges"
	"dogker/andrenk/billing-service/internal/rest/deposits"
	"dogker/andrenk/billing-service/internal/rest/mutations"
	"dogker/andrenk/billing-service/internal/rest/payment"
	"fmt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	//GORM Setup
	//TODO: 2
	dsn := "host=dogker-postgres user=admin password=admin dbname=dogker port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	db.AutoMigrate(&deposits.Deposit{}, &mutations.Mutation{}, &charges.Charge{})

	//Public Key Setup
	//TODO: 3
	publicKeyAuthServer := "-----BEGIN PUBLIC KEY-----\nMIGbMBAGByqGSM49AgEGBSuBBAAjA4GGAAQAodxwFdiFKWTG/ZU7vXPdk8ox+nNU\n1JmxsmI8i8tYrYf6QxmwBz13jS/PZsb8dJbMFY3YTMMih6SKz7e+cQ68IbgA7BnY\n5fYFQET4SNHVX/zaH6J70ERJLsRrarmWSXsNbMbnqXlIkoorYXeAn9vsLbr/RPw9\nDYaoq4JrQ+OGsc4LHMw=\n-----END PUBLIC KEY-----\n"

	err = auth.InitPublicKey(publicKeyAuthServer)
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
	mutationsHandler := mutations.NewMutationsHandler(mutationService)

	//Deposits Setup
	depositRepository := deposits.NewRepository(db)
	depositService := deposits.NewService(depositRepository)
	depositsHandler := deposits.NewDepositsHandler(depositService, paymentService, authService, mutationService)

	//Charges Setup
	chargeRepository := charges.NewRepository(db)
	chargeService := charges.NewService(chargeRepository, mutationService)
	chargeHandler := charges.NewChargeHandler(chargeService, mutationService)

	//RabbitMQ Setup
	rmq := rabbitmq.NewRabbitMQ("amqp://guest:guest@rabbitmq:5672/")
	msgs, err := rmq.Channel.Consume(
		"monitor-billing",  // queue
		"billing-consumer", // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	if err != nil {
		_ = fmt.Errorf("err: rmq.Channel.Consume: %w", err)
	}

	go func() {
		for msg := range msgs {
			allUsersMetrics, err := rabbitmq.DecodeAllUserMetricsMessage(msg.Body)
			if err != nil {
				_ = fmt.Errorf("decodeAllUserMetricsMessage: %w", err)
			}

			chargeService.ChargeInBatch(allUsersMetrics)
		}
	}()

	//REST Setup
	router := gin.Default()
	api := router.Group("/api/v1")

	router.RedirectTrailingSlash = false
	corsSeting := cors.Default()

	router.Use(corsSeting)

	depositsRoute := api.Group("/deposits")
	{
		depositsRoute.POST("", rest.AuthMiddleware(), rest.CORSMiddleware(), depositsHandler.InitiateDeposit)
		depositsRoute.GET("/:id", rest.AuthMiddleware(), rest.CORSMiddleware(), depositsHandler.GetDepositByID)
	}

	chargesRoute := api.Group("/charges")
	{
		chargesRoute.POST("", rest.AuthMiddleware(), rest.CORSMiddleware(), chargeHandler.InitiateCharge)
		chargesRoute.GET("/:id", rest.AuthMiddleware(), rest.CORSMiddleware(), chargeHandler.GetChargeByID)
	}

	mutationsRoute := api.Group("/mutations")
	{
		mutationsRoute.GET("", rest.AuthMiddleware(), rest.CORSMiddleware(), mutationsHandler.GetMutationsByUserID)
		mutationsRoute.GET("/:id", rest.AuthMiddleware(), rest.CORSMiddleware(), mutationsHandler.GetMutationByID)
	}

	notificationRoute := api.Group("/notification")
	{
		notificationRoute.POST("", depositsHandler.HandleNotificationPayment)
	}

	router.Run(":6969")
}
