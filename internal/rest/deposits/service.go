package deposits

import (
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
)

type Service interface {
	InitiateDeposit(deposit DepositInput) (*snap.Response, error)
	AddDeposit(deposit DepositInput) (Deposit, error)
	GetDepositByID(id string) (Deposit, error)
	UpdateDeposit(deposit DepositUpdateInput) (Deposit, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) InitiateDeposit(deposit DepositInput) (*snap.Response, error) {
	var snapInstance = snap.Client{}

	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")
	if serverKey == "" {
		return nil, fmt.Errorf("server key is not found")
	}

	snapInstance.New(serverKey, midtrans.Sandbox)

	req := snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  uuid.NewString(),
			GrossAmt: int64(deposit.Amount),
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: "John",
			LName: "Doe",
			Email: "johndoe@gmail.com",
			Phone: "081234567890",
		},
	}

	snapResponse, _ := s.repository.Initiate(snapInstance, &req)

	return snapResponse, nil
}

func (s *service) AddDeposit(deposit DepositInput) (Deposit, error) {
	depositObj := Deposit{}

	depositObj.ID = uuid.New().String()
	depositObj.UserID = deposit.UserID
	depositObj.Amount = deposit.Amount

	return s.repository.Save(depositObj)
}

func (s *service) GetDepositByID(id string) (Deposit, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateDeposit(deposit DepositUpdateInput) (Deposit, error) {
	depositObj, err := s.repository.FindByID(deposit.ID)

	if err != nil {
		return Deposit{}, err
	}

	depositObj.UserID = deposit.UserID
	depositObj.Amount = deposit.Amount
	depositObj.Status = deposit.Status

	return s.repository.Update(depositObj)
}
