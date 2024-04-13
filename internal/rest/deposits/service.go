package deposits

import "github.com/google/uuid"

type Service interface {
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
