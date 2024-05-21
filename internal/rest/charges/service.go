package charges

import "time"

type Service interface {
	CreateCharge(charge ChargeInput) (Charge, error)
	GetChargeByID(id string) (Charge, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) CreateCharge(charge ChargeInput) (Charge, error) {
	chargeObj := Charge{
		ID:                   charge.ID,
		UserID:               charge.UserID,
		ContainerID:          charge.ContainerID,
		TotalCpuUsage:        charge.TotalCpuUsage,
		TotalMemoryUsage:     charge.TotalMemoryUsage,
		TotalNetIngressUsage: charge.TotalNetIngressUsage,
		TotalNetEgressUsage:  charge.TotalNetEgressUsage,
		Timestamp:            time.Now(),
		TotalCost:            charge.TotalCost,
	}

	return s.repository.Save(chargeObj)
}

func (s *service) GetChargeByID(id string) (Charge, error) {
	return s.repository.FindByID(id)
}
