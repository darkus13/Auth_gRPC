package auth

import (
	"github.com/darkus13/Auth_gRPC/internal/service"
	decs "github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

type Implementation struct {
	decs.UnimplementedUserV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{
		authService: authService,
	}
}
