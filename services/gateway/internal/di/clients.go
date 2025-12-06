package di

import (
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSC struct {
	Auth authv1.AuthServiceClient
}

func NewGRPCServiceClients(authAddr string) (GRPCSC, error) {
	authCli, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return GRPCSC{}, err
	}

	return GRPCSC{
		Auth: authv1.NewAuthServiceClient(authCli),
	}, nil

}
