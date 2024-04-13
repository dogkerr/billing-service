package main

import (
	"dogker/andrenk/billing-service/internal/rest"
	"dogker/andrenk/billing-service/internal/rest/deposits"
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

	db.AutoMigrate(&deposits.Deposit{})

	depositRepository := deposits.NewRepository(db)
	depositService := deposits.NewService(depositRepository)
	depositsHandler := deposits.NewDepositsHandler(depositService)

	router := gin.Default()
	api := router.Group("/api/v1")

	router.Use(rest.AuthMiddleware())

	depositsRoute := api.Group("/deposits")
	{
		depositsRoute.GET("/deposits", depositsHandler.GetDepositByID)
		depositsRoute.POST("/deposits", depositsHandler.AddDeposit)
		depositsRoute.PUT("/deposits", depositsHandler.UpdateDeposit)
	}

	router.Run(":8080")
}
