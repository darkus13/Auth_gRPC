package app

import (
	"context"
	"log"

	"github.com/darkus13/Auth_gRPC/internal/api/auth"
	"github.com/darkus13/Auth_gRPC/internal/client/db"
	"github.com/darkus13/Auth_gRPC/internal/client/db/pg"
	"github.com/darkus13/Auth_gRPC/internal/closer"
	"github.com/darkus13/Auth_gRPC/internal/config"
	"github.com/darkus13/Auth_gRPC/internal/repository"
	authRepository "github.com/darkus13/Auth_gRPC/internal/repository/auth"
	"github.com/darkus13/Auth_gRPC/internal/service"
	authService "github.com/darkus13/Auth_gRPC/internal/service/auth"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient       db.Client
	authService    service.AuthService
	authRepository repository.AuthRepository

	authImplm *auth.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := config.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}

func (s *serviceProvider) GRPCConfig() config.GRPCConfig {
	if s.grpcConfig == nil {
		cfg, err := config.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to get gRPC config: %s", err.Error())
		}

		s.grpcConfig = cfg
	}

	return s.grpcConfig
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) AuthRepository(ctx context.Context) repository.AuthRepository {
	if s.authRepository == nil {
		s.authRepository = authRepository.NewRepository(s.DBClient(ctx).DB())
	}

	return s.authRepository
}

func (s *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewService(s.AuthRepository(ctx))
	}

	return s.authService
}

func (s *serviceProvider) AuthImplm(ctx context.Context) *auth.Implementation {
	if s.authImplm == nil {
		s.authImplm = auth.NewImplementation(s.AuthService(ctx))
	}

	return s.authImplm
}
