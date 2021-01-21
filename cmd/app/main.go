package main

import (
	"context"
	"github.com/zhashkevych/courses-backend/internal/config"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
	"github.com/zhashkevych/courses-backend/pkg/logger"
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

	if err := logger.Init(); err != nil {
		panic(err)
	}

	logger.Infof("%+v\n", cfg)

	logger.Info(db.CreateCollection(context.Background(), "users"))
}
