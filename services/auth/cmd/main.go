package main

import (
	"context"
	"errors"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"

	authv1 "github.com/ummuys/reportify/api/pb/auth/service/v1"
	"github.com/ummuys/reportify/pkg/config"
	"github.com/ummuys/reportify/pkg/errs"
	"github.com/ummuys/reportify/pkg/logger"
	pkg "github.com/ummuys/reportify/pkg/tm"
	"github.com/ummuys/reportify/services/auth/internal/adapter"
	"github.com/ummuys/reportify/services/auth/internal/dto"
	"github.com/ummuys/reportify/services/auth/internal/repository"
	"github.com/ummuys/reportify/services/auth/internal/secure"
	"github.com/ummuys/reportify/services/auth/internal/service"
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
	tm, err := pkg.NewTokenManager()
	if err != nil {
		logs.Fatal().Err(err).Msg("token-manager")
	}

	db, err := repository.NewAuthDB(ctx, logs)
	if err != nil {
		logs.Fatal().Err(err).Msg("auth-db")
	}
	defer db.Close()

	srv := grpc.NewServer()
	svc := service.NewAuthService(ph, tm, db, logs)

	// CREATE USER
	if err := svc.CreateBaseAdmin(ctx, dto.CreateUserParams{
		Username: cfg.AdmUsername,
		Password: cfg.AdmPassword,
		Role:     "admin",
	}); err != nil && !errors.Is(err, errs.ErrDuplicate) {
		logs.Fatal().Err(err).Msg("create admin user")
	}

	adp := adapter.NewAuthAdapter(svc, logs)
	authv1.RegisterAuthServiceServer(srv, adp)

	wg := sync.WaitGroup{}
	SDChan := make(chan errs.SDMsg, 2)

	wg.Go(func() {
		<-ctx.Done()
		srv.GracefulStop()
	})

	wg.Go(func() {
		logs.Info().Msg("run the grpc-server")
		if err := srv.Serve(lis); err != nil {
			SDChan <- errs.SDMsg{
				Err:  err,
				From: "grpc-server",
			}
		}
	})

	wg.Wait()
	close(SDChan)
	errs.ShutdownStatus(logs, SDChan)
}
