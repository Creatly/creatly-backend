package main

import (
	"context"
	"fmt"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
)

func main() {
	mongoClient := mongodb.NewClient("mongodb://mongodb:27017", "admin", "qwerty")
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("coursePlatform")

	fmt.Println(db.CreateCollection(context.Background(), "users"))
}
