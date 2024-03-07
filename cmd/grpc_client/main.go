package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/brianvoe/gofakeit"

	"github.com/darkus13/Auth_gRPC/pkg/user_api_v1"
)

const (
	address = "localhost:50051"
	userID  = 2
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	c := user_api_v1.NewUserV1Client(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	testPassword := gofakeit.Password(true, false, true, true, false, 6)

	r, err := c.Create(ctx, &user_api_v1.CreateRequest{
		Name:            gofakeit.Name(),
		Email:           gofakeit.Email(),
		Password:        testPassword,
		PasswordConfirm: testPassword,
		Role:            user_api_v1.Role(gofakeit.Number(0, 1)),
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	log.Printf("user with id: %d was created.", r.GetId())
}
