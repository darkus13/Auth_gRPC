package auth

import (
	"context"
	"log"

	"github.com/darkus13/Auth_gRPC/internal/converter"
	decs "github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

func (i *Implementation) Create(ctx context.Context, req *decs.CreateRequest) (*decs.CreateResponse, error) {
	id, err := i.authService.Create(ctx, converter.ToServiceFromAuth(req))
	if err != nil {
		return nil, err
	}

	log.Printf("inserted auth with id: %v", id)

	return &decs.CreateResponse{
		Id: id,
	}, nil
}
