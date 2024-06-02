package charges

import (
	"context"
	"dogker/andrenk/billing-service/internal/grpc"
	"dogker/andrenk/billing-service/internal/rabbitmq"
	"dogker/andrenk/billing-service/internal/rest/mutations"
	"dogker/andrenk/billing-service/protos"
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service interface {
	CreateCharge(charge ChargeInput) (Charge, error)
	ChargeInBatch(charge []rabbitmq.UserMetricsMessage) error
	GetChargeByID(id string) (Charge, error)
}

type service struct {
	repository      Repository
	mutationService mutations.Service
}

func NewService(repository Repository, mutationService mutations.Service) *service {
	return &service{repository, mutationService}
}

func (s *service) CreateCharge(charge ChargeInput) (Charge, error) {
	//Get mutation list
	fmt.Println("Getting mutations...")
	mutationsList, err := s.mutationService.GetMutationsByUserID(charge.UserID)
	if err != nil {
		_ = fmt.Errorf("error getting mutations: %v", err)
		return Charge{}, err
	}

	//Calculate cost
	totalCost := 0.0175*charge.TotalCpuUsage + 0.08*charge.TotalMemoryUsage + 0.00175*charge.TotalNetIngressUsage + 0.00175*charge.TotalNetEgressUsage
	fmt.Println("Total cost: ", totalCost)

	//Check if user has enough balance, if not stop container
	if (len(mutationsList) == 0) || (mutationsList[0].Balance < totalCost) {
		//Stop container
		fmt.Printf("stopping user %s contianer...", charge.UserID)
		conn, err := grpc.GetGRPCClient("container-service:6666")
		if err != nil {
			err = fmt.Errorf("error stopping container: %v", err)
			return Charge{}, err
		}

		defer conn.Close()

		containerClient := protos.NewContainerGRPCServiceClient(conn)
		_, err = containerClient.StopContainerCreditLimit(context.Background(), &protos.StopUserContainerCreditLimitReq{
			UserID: charge.UserID,
		})

		if err != nil {
			err = fmt.Errorf("error stopping container: %v", err)
			return Charge{}, err
		}

		err = fmt.Errorf("user does not have enough balance")
		return Charge{}, err
	}

	//Charge user
	fmt.Println("Charging user...")
	chargeObj := Charge{
		ID:                   charge.ID,
		UserID:               charge.UserID,
		ContainerID:          charge.ContainerID,
		TotalCpuUsage:        charge.TotalCpuUsage,
		TotalMemoryUsage:     charge.TotalMemoryUsage,
		TotalNetIngressUsage: charge.TotalNetIngressUsage,
		TotalNetEgressUsage:  charge.TotalNetEgressUsage,
		Timestamp:            time.Now(),
		TotalCost:            totalCost,
	}
	chargeResult, error := s.repository.Save(chargeObj)

	//Create mutation
	fmt.Println("Creating mutation...")
	mutationInput := mutations.MutationInput{
		ID:       uuid.NewString(),
		UserID:   chargeObj.UserID,
		Mutation: chargeObj.TotalCost,
		Type:     "charge",
		ChargeID: chargeObj.ID,
	}
	s.mutationService.CreateMutation(mutationInput)

	//Call mailing service
	fmt.Println("Sending email...")
	conn, err := grpc.GetGRPCClient("mailing-service:9897")
	if err != nil {
		err = fmt.Errorf("error stopping container: %v", err)
		return Charge{}, err
	}
	defer conn.Close()

	emailClient := protos.NewEmailServiceClient(conn)
	_, err = emailClient.SendBillingEmail(context.Background(), &protos.BillingEmailRequest{
		Id:                   charge.ID,
		UserId:               charge.UserID,
		ContainerId:          charge.ContainerID,
		TotalCpuUsage:        charge.TotalCpuUsage,
		TotalMemoryUsage:     charge.TotalMemoryUsage,
		TotalNetIngressUsage: charge.TotalNetIngressUsage,
		TotalNetEgressUsage:  charge.TotalNetEgressUsage,
		TotalCost:            totalCost,
		Timestamp:            timestamppb.Now(),
	})

	if err != nil {
		err = fmt.Errorf("error while calling mailing service: %v", err)
		return Charge{}, err
	}

	return chargeResult, error
}

func (s *service) ChargeInBatch(metrics []rabbitmq.UserMetricsMessage) error {
	var err error

	for _, metric := range metrics {
		chargeInput := ChargeInput{
			ID:                   uuid.NewString(),
			UserID:               metric.UserID,
			ContainerID:          metric.ContainerID,
			TotalCpuUsage:        float32(metric.CpuUsage),
			TotalMemoryUsage:     float32(metric.MemoryUsage),
			TotalNetIngressUsage: float32(metric.NetworkIngressUsage),
			TotalNetEgressUsage:  float32(metric.NetworkEgressUsage),
		}

		_, _err := s.CreateCharge(chargeInput)
		if _err != nil {
			err = _err
		}
	}

	return err
}

func (s *service) GetChargeByID(id string) (Charge, error) {
	return s.repository.FindByID(id)
}
