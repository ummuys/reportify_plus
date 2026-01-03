package di

import (
	authv1 "github.com/ummuys/reportify/api/pb/auth/service/v1"
	dsv1 "github.com/ummuys/reportify/api/pb/datasource/service/v1"
	rsv1 "github.com/ummuys/reportify/api/pb/report/service/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCSC struct {
	AuthService       authv1.AuthServiceClient
	ReportService     rsv1.ReportServiceClient
	DatasourceService dsv1.DatasourceServiceClient
}

func NewGRPCServiceClients(authAddr, reportAddr string) (GRPCSC, error) {
	authCli, err := grpc.NewClient(
		authAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return GRPCSC{}, err
	}

	reportSvcCli, err := grpc.NewClient(
		reportAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return GRPCSC{}, err
	}

	datasourceCli, err := grpc.NewClient(
		reportAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return GRPCSC{}, err
	}

	return GRPCSC{
		AuthService:       authv1.NewAuthServiceClient(authCli),
		ReportService:     rsv1.NewReportServiceClient(reportSvcCli),
		DatasourceService: dsv1.NewDatasourceServiceClient(datasourceCli),
	}, nil
}
