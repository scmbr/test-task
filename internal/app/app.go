package app

import (
	"github.com/scmbr/test-task/internal/config"
	delivery "github.com/scmbr/test-task/internal/delivery/http"
	"github.com/scmbr/test-task/internal/repository"
	"github.com/scmbr/test-task/internal/server"
	"github.com/scmbr/test-task/internal/service"
	"github.com/scmbr/test-task/pkg/auth"
	"github.com/scmbr/test-task/pkg/database"
	"github.com/scmbr/test-task/pkg/hasher"
	"github.com/sirupsen/logrus"
)

func Run(configsDir string) {
	cfg, err := config.Init(configsDir)
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
		return
	}
	db, err := database.NewPostgresDB(database.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		Username: cfg.Postgres.Username,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.Name,
		SSLMode:  cfg.Postgres.SSLMode,
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db:%s", err.Error())
	}
	hasher := hasher.NewHasher(cfg.Hasher.Cost)
	repositories := repository.NewRepository(db)
	tokenManager, err := auth.NewManager(cfg.Auth.SigningKey)
	if err != nil {
		logrus.Fatalf("error initializing token manager: %s", err.Error())
		return
	}
	services := service.NewServices(service.Deps{
		Repos:           repositories,
		Hasher:          *hasher,
		TokenManager:    tokenManager,
		AccessTokenTTL:  cfg.Auth.AccessTokenTTL,
		RefreshTokenTTL: cfg.Auth.RefreshTokenTTL,
	})
	srv := new(server.Server)
	handlers := delivery.NewHandler(services)
	if err := srv.Run(cfg.HTTP.Port, cfg.HTTP.Host, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while starting http server: %s", err.Error())
	}
}
