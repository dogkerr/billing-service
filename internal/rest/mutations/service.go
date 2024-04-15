package mutations

import "time"

type Service interface {
	GetMutations() ([]Mutation, error)
	GetMutationByID(id string) (Mutation, error)
	CreateMutation(mutation MutationInput) (Mutation, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetMutations() ([]Mutation, error) {
	return s.repository.GetMutations()
}

func (s *service) GetMutationByID(id string) (Mutation, error) {
	return s.repository.GetMutationByID(id)
}

func (s *service) CreateMutation(mutation MutationInput) (Mutation, error) {
	//Get previous mutation
	mutations, err := s.repository.GetMutationByUserID(mutation.UserID)
	if err != nil {
		return Mutation{}, err
	}

	//Calculate current balance
	currentBalance := float64(0)
	if len(mutations) > 0 {
		currentBalance = mutations[0].Balance
		if mutation.Type == "charge" {
			currentBalance -= mutation.Mutation
		} else if mutation.Type == "deposit" {
			currentBalance += mutation.Mutation
		}
	} else {
		if mutation.Type == "charge" {
			currentBalance -= mutation.Mutation
		} else if mutation.Type == "deposit" {
			currentBalance += mutation.Mutation
		}
	}

	mutationObj := Mutation{
		ID:        mutation.ID,
		UserID:    mutation.UserID,
		Mutation:  mutation.Mutation,
		Balance:   currentBalance,
		Type:      mutation.Type,
		DepositID: mutation.DepositID,
		ChargeID:  mutation.ChargeID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	return s.repository.CreateMutation(mutationObj)
}
