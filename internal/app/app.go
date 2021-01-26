package app

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zhashkevych/courses-backend/internal/config"
	"github.com/zhashkevych/courses-backend/internal/delivery/http"
	"github.com/zhashkevych/courses-backend/internal/repository"
	"github.com/zhashkevych/courses-backend/internal/server"
	"github.com/zhashkevych/courses-backend/internal/service"
	"github.com/zhashkevych/courses-backend/pkg/auth"
	"github.com/zhashkevych/courses-backend/pkg/cache"
	"github.com/zhashkevych/courses-backend/pkg/database/mongodb"
	"github.com/zhashkevych/courses-backend/pkg/email/sendpulse"
	"github.com/zhashkevych/courses-backend/pkg/hash"
	"github.com/zhashkevych/courses-backend/pkg/logger"
	"github.com/zhashkevych/courses-backend/pkg/payment"
	"os"
	"os/signal"
	"syscall"
)

// @title Course Platform API
// @version 1.0
// @description API Server for Course Platform

// TODO host based on env
// @host localhost:8000
// @BasePath /api/v1/

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
		logger.Error(err)
		return
	}

	// Dependencies
	mongoClient := mongodb.NewClient(cfg.Mongo.URI, cfg.Mongo.User, cfg.Mongo.Password)
	db := mongoClient.Database(cfg.Mongo.Name)

	memCache := cache.NewMemoryCache(int64(cfg.CacheTTL))
	hasher := hash.NewSHA1Hasher(cfg.Auth.PasswordSalt)
	emailProvider := sendpulse.NewClient(cfg.Email.ClientID, cfg.Email.ClientSecret, memCache)
	tokenManager, err := auth.NewManager(cfg.Auth.JWT.SigningKey)
	if err != nil {
		logger.Error(err)
		return
	}

	// Services, Repos & API Handlers
	repos := repository.NewRepositories(db)
	services := service.NewServices(repos, memCache, hasher, tokenManager,
		emailProvider, cfg.Email.ListID, payment.MockProvider{}, cfg.Auth.JWT.AccessTokenTTL, cfg.Auth.JWT.RefreshTokenTTL)
	handlers := http.NewHandler(services.Schools, services.Students, services.Courses, tokenManager)

	// HTTP Server
	srv := server.NewServer(cfg, handlers.Init(cfg.HTTP.Host, cfg.HTTP.Port))
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
