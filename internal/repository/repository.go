package repository

import (
	"context"

	desc "github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

type AuthRepository interface {
	Create(ctx context.Context, info *desc.CreateResponse) (int64, error)
	Get(ctx context.Context, id int64) (*desc.Auth, error)
}
