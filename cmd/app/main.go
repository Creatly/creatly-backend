package main

import (
	"context"
	"fmt"
	"github.com/zhashkevych/courses-backend/internal/config"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
)

const configPath = "configs/main"

func main() {
	mongoClient := mongodb.NewClient("mongodb://mongodb:27017", "admin", "qwerty")
	defer mongoClient.Disconnect(context.Background())

	db := mongoClient.Database("coursePlatform")

	cfg, err := config.Init(configPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)

	fmt.Println(db.CreateCollection(context.Background(), "users"))
}
