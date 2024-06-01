package mutations

import "time"

type Service interface {
	GetMutationsByUserID(id string) ([]Mutation, error)
	GetMutationByID(id string) (Mutation, error)
	CreateMutation(mutation MutationInput) (Mutation, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) GetMutationByID(id string) (Mutation, error) {
	return s.repository.FindByID(id)
}

func (s *service) GetMutationsByUserID(id string) ([]Mutation, error) {
	return s.repository.FindByUserID(id)
}

func (s *service) CreateMutation(mutation MutationInput) (Mutation, error) {
	//Get previous mutation
	mutations, err := s.repository.FindByUserID(mutation.UserID)
	if err != nil {
		return Mutation{}, err
	}

	//Calculate current balance
	var currentBalance float32
	if mutation.Type == "charge" && len(mutations) > 0 && mutations[0].Balance > mutation.Mutation {
		currentBalance = float32(mutations[0].Balance)
		currentBalance -= mutation.Mutation
	} else if mutation.Type == "deposit" {
		if len(mutations) < 1 {
			currentBalance = float32(0)
			currentBalance += mutation.Mutation
		} else {
			currentBalance = float32(mutations[0].Balance)
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

	return s.repository.Save(mutationObj)
}
