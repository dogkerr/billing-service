package auth

import (
	"context"
	"dogker/andrenk/billing-service/internal/grpc"
	"dogker/andrenk/billing-service/protos"
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
	conn, err := grpc.GetGRPCClient("103.175.219.0:4000")
	if err != nil {
		return nil, err
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
