package auth

import (
	"context"
	"dogker/andrenk/billing-service/protos"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type service struct{}

type Service interface {
	GetUserById(userID string) (*protos.User, error)
}

func NewService() *service {
	return &service{}
}

func (s *service) GetUserById(userID string) (*protos.User, error) {
	//users gRPC Client Setup
	conn, err := grpc.Dial("10.66.66.1:4001", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}
	defer conn.Close()

	usersClient := protos.NewUsersServiceClient(conn)

	req := &protos.GetUserRequest{
		Id: userID,
	}

	//Get User by ID
	user, err := usersClient.GetUserById(context.Background(), req)

	if err != nil {
		return nil, err
	}

	return user, nil
}
