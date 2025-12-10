package di

import (
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	reportv1 "github.com/ummuys/reportify/api/pb/report/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSC struct {
	Auth   authv1.AuthServiceClient
	Report reportv1.ReportServiceClient
}

func NewGRPCServiceClients(authAddr, reportAddr string) (GRPCSC, error) {
	authCli, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return GRPCSC{}, err
	}

	reportCli, err := grpc.NewClient(
		reportAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return GRPCSC{}, err
	}
	return GRPCSC{
		Report: reportv1.NewReportServiceClient(reportCli),
		Auth:   authv1.NewAuthServiceClient(authCli),
	}, nil

}
