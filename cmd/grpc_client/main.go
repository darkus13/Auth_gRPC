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
	defer conn.Close()

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

	// r, err := c.Get(ctx, &desc.GetRequest{Id: userID})
	// if err != nil {
	// 	log.Fatalf("failed to get user by id: %v", err)
	// }

	// log.Printf("User info:\nID: %d\nName: %s\nEmail: %s\nRole: %s\nCreated at: %v\nUpdated at: %v",
	// 	r.GetId(), r.GetName(), r.GetEmail(), r.GetRole(),
	// 	r.GetCreatedAt(), r.GetUpdatedAt())

	// _, err = c.Delete(ctx, &desc.DeleteRequest{Id: userID})
	// if err != nil {
	// 	log.Fatalf("failed to delete user: %v", err)
	// }

	// _, err = c.Update(ctx, &desc.UpdateRequest{
	// 	Id:    userID,
	// 	Name:  wrapperspb.String(gofakeit.Name()),
	// 	Email: wrapperspb.String(gofakeit.Email()),
	// 	Role:  desc.UserRole(1),
	// })
	// if err != nil {
	// 	log.Fatalf("failed to update user info: %v", err)
	// }
}
