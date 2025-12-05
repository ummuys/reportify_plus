package main

import (
	"context"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/pkg/logger"
	"github.com/ummuys/reportify/services/auth/internal/adapter"
	"github.com/ummuys/reportify/services/auth/internal/auth"
	"google.golang.org/grpc"
)

func main() {
	mainCtx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logs, err := logger.InitLogger("auth")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logs.Fatal().Err(err).Msg("")
	}

	srv := grpc.NewServer()
	svc := auth.NewAuthService()
	adp := adapter.NewAuthAdapter(svc)
	authv1.RegisterAuthServiceServer(srv, adp)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		<-mainCtx.Done()
		srv.GracefulStop()
	})

	wg.Go(func() {
		if err := srv.Serve(lis); err != nil {
			logs.Error().Err(err).Msg("")
		}
	})

	wg.Wait()

}
