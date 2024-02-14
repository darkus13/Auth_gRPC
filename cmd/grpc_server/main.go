package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/brianvoe/gofakeit/v6"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"grpc/pkg/user_api_v1"
)

const grpcPort = 50051

// Структура сервера
type server struct {
	user_api_v1.UnimplementedUserV1Server
}

func (s *server) Get(ctx context.Context, req *user_api_v1.GetRequest) (*user_api_v1.GetResponse, error) {
	log.Printf("Note id: %d", req.GetId())
	return &user_api_v1.GetResponse{
		Id:        req.GetId(),
		Name:      gofakeit.Name(),
		Email:     gofakeit.Email(),
		Role:      user_api_v1.Role_USER,
		CreatedAt: timestamppb.New(gofakeit.Date()),
		UpdatedAt: timestamppb.New(gofakeit.Date()),
	}, nil
}

func (s *server) Create(ctx context.Context, req *user_api_v1.CreateRequest) (*user_api_v1.CreateResponse, error) {
	return &user_api_v1.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Update(ctx context.Context, req *user_api_v1.UpdateRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

// Delete Cоздадим удаление пользователя.
func (s *server) Delete(ctx context.Context, req *user_api_v1.DeleteRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	user_api_v1.RegisterUserV1Server(s, &server{})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
