package auth

import (
	"github.com/darkus13/Auth_gRPC/internal/repository"
	"github.com/darkus13/Auth_gRPC/internal/service"
)

type serv struct {
	authRepository repository.AuthRepository
}

func NewService(
	authRepository repository.AuthRepository) service.AuthService {
	return &serv{
		authRepository: authRepository,
	}
}
