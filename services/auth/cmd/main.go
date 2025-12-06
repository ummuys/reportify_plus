package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	"github.com/joho/godotenv"
	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/pkg/logger"
	"github.com/ummuys/reportify/services/auth/internal/adapter"
	"github.com/ummuys/reportify/services/auth/internal/auth"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/repository"
	"github.com/ummuys/reportify/services/auth/internal/secure"
	"google.golang.org/grpc"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	_ = godotenv.Load("../../../.env")

	logs, err := logger.InitLogger("auth", "LOG_LEVEL_AUTH")
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logs.Fatal().Err(err).Msg("listener")
	}

	ph := secure.NewPasswordHasher()
	tm, err := secure.NewTokenManager()

	if err != nil {
		logs.Fatal().Err(err).Msg("token-manager")
	}

	db, err := repository.NewAuthDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("auth-db")
	}

	srv := grpc.NewServer()
	svc := auth.NewAuthService(ph, tm, db, logs)

	// CREATE USER
	if out, err := svc.CreateUser(ctx, dto.CreateUserParams{
		Username: "admin",
		Password: "admin",
		Role:     "admin",
	}); err != nil {
		logs.Error().Err(err).Msg("create user")
	} else {
		fmt.Println(out.UserID)
	}
	adp := adapter.NewAuthAdapter(svc, logs)
	authv1.RegisterAuthServiceServer(srv, adp)

	wg := sync.WaitGroup{}
	wg.Go(func() {
		<-ctx.Done()
		srv.GracefulStop()
	})

	wg.Go(func() {
		logs.Info().Msg("run the grpc-server")
		if err := srv.Serve(lis); err != nil {
			logs.Error().Err(err).Msg("")
		}
	})

	wg.Wait()
	logs.Info().Msg("graceful shutdown")
}
