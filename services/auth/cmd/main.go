package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	authv1 "github.com/ummuys/reportify/api/pb/auth/v1"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/errs"
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

	logs, err := logger.InitLogger("auth", "LOG_LEVEL_AUTH")
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := config.ParseAuthServiceConfig()
	if err != nil {
		logs.Fatal().Err(err).Msg("config")
	}

	lis, err := net.Listen(cfg.Network, cfg.Port)
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
	if _, err := svc.CreateUser(ctx, dto.CreateUserParams{
		Username: cfg.AdmUsername,
		Password: cfg.AdmPassword,
		Role:     "admin",
	}); err != nil && !errors.Is(err, errs.ErrDuplicate) {
		logs.Fatal().Err(err).Msg("create admin user")
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
			logs.Error().Err(err).Msg("grpc-server")
		}
	})

	wg.Wait()
	logs.Info().Msg("graceful shutdown")
}
