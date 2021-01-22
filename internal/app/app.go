package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/courses-backend/internal/config"
	"github.com/zhashkevych/courses-backend/internal/delivery/http"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/internal/server"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"os"
	"os/signal"
	"syscall"
)

// @title Course Platform API
// @version 1.0
// @description API Server for Course Platform

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey AdminAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey StudentsAuth
// @in header
// @name Authorization

// Run initializes whole application
func Run(configPath string) {
	cfg, err := config.Init(configPath)
	if err != nil {
		panic(err)
	}

	logger.Infof("%+v", *cfg)

	mongoClient := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)

	db := mongoClient.Database(cfg.Mongo.Name)
	memCache := cache.NewMemoryCache()
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)

	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, memCache, hasher)

	handlers := http.NewHandler(services.Schools, services.Students)

	srv := server.NewServer(cfg, handlers.Init())
	go func() {
		if err := srv.Run(); err != nil {
			logrus.Errorf("error occurred while running http server: %s\n", err.Error())
		}
	}()

	logger.Info("Server started")

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	if err := mongoClient.Disconnect(context.Background()); err != nil {
		logger.Error(err.Error())
	}
}
