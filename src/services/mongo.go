package services

import (
	"context"
	"fmt"
	"log"

	"github.com/weeber-id/desatanjungbunga-backend/src/variables"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	// Client for mongodb connection
	Client *mongo.Client

	// DB for mongoDB Database pointer
	DB *mongo.Database
)

// InitializationMongo for mongodb connection initialization
func InitializationMongo(ctx context.Context) *mongo.Client {
	config := variables.MongoConfig

	URI := fmt.Sprintf(
		"%s://%s:%s@%s/%s?retryWrites=true&w=majority",
		config.Connector,
		config.User,
		config.Password,
		config.Host,
		config.Database,
	)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatalf("Error connecting mongoDB: %v", err)
	}
	Client = client

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Error ping mongoDB: %v", err)
	}

	DB = client.Database(config.Database)

	return client
}
