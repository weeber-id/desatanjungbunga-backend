package main

import (
	"context"
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/weeber-id/desatanjungbunga-backend/src/services"
	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
)

func main() {
	// Environment section
	godotenv.Load("./devel.env")
	variables.InitializationVariable()

	// MongoDB section
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	client := services.InitializationMongo(ctx)
	defer client.Disconnect(ctx)

	fmt.Println("finish")
}
