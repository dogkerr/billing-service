package payment

import (
	"dogker/andrenk/billing-service/protos"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type service struct{}

type Service interface {
	GetPayment(input TransactionInput, user *protos.User) (Transaction, error)
	HandleNotificationPayment(input TransactionNotificationInput) bool
}

func NewService() *service {
	return &service{}
}

func (s *service) GetPayment(input TransactionInput, user *protos.User) (Transaction, error) {
	fmt.Println("GetPayment")

	//Midtrans Setup
	env, err := godotenv.Read(".env")
	if err != nil {
		return Transaction{}, err
	}

	serverKey := env["MIDTRANS_SERVER_KEY"]
	if serverKey == "" {
		return Transaction{}, fmt.Errorf("MIDTRANS_SERVER_KEY is not found")
	}
	fmt.Println("Server Key: ", serverKey)

	var snapClient snap.Client
	snapClient.New(serverKey, midtrans.Sandbox)

	orderID := uuid.NewString()
	req := snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  orderID,
			GrossAmt: int64(input.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: user.Fullname,
			Email: user.Email,
		},
	}

	snapResponse, _ := snapClient.CreateTransaction(&req)
	transaction := Transaction{
		ID:          orderID,
		UserID:      user.Id,
		Amount:      input.Amount,
		Token:       snapResponse.Token,
		RedirectURL: snapResponse.RedirectURL,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return transaction, nil
}

func (s *service) HandleNotificationPayment(input TransactionNotificationInput) bool {
	if input.PaymentType == "credit_card" && input.TransactionStatus == "capture" && input.FraudStatus == "accept" {
		return true
	} else if input.TransactionStatus == "settlement" {
		return true
	} else if input.TransactionStatus == "deny" || input.TransactionStatus == "expire" || input.TransactionStatus == "cancel" {
		return false
	}

	return false
}
