package mongodb

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const timeout = 10 * time.Second

// NewClient established connection to a mongoDb instance using provided URI and auth credentials
func NewClient(uri, username, password string) *mongo.Client {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: username, Password: password}))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
