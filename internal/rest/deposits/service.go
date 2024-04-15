package deposits

import "time"

type Service interface {
	AddDeposit(deposit DepositInput) (Deposit, error)
	GetDepositByID(id string) (Deposit, error)
	UpdateDepositStatus(id string, status string) (Deposit, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) AddDeposit(deposit DepositInput) (Deposit, error) {
	depositObj := Deposit{
		ID:        deposit.ID,
		UserID:    deposit.UserID,
		Amount:    deposit.Amount,
		Status:    deposit.Status,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repository.Save(depositObj)
}

func (s *service) GetDepositByID(id string) (Deposit, error) {
	return s.repository.FindByID(id)
}

func (s *service) UpdateDepositStatus(id string, status string) (Deposit, error) {
	depositObj, err := s.repository.FindByID(id)
	if err != nil {
		return Deposit{}, err
	}

	depositObj.Status = status
	depositObj.UpdatedAt = time.Now()

	return s.repository.Update(depositObj)
}
