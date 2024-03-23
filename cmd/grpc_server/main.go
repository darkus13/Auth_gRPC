package main

import (
	"context"
	"log"

	"github.com/darkus13/Auth_gRPC/internal/app"
)

func main() {
	ctx := context.Background()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Fatalf("failed to init app: #{err.Error()}")
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: #{err.Error()}")
	}
}
