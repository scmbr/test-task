package app

import (
	"github.com/scmbr/test-task/internal/config"
	delivery "github.com/scmbr/test-task/internal/delivery/http"
	"github.com/scmbr/test-task/internal/server"
	"github.com/scmbr/test-task/pkg/database"
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
	srv := new(server.Server)
	handlers := delivery.NewHandler()
	if err := srv.Run(cfg.HTTP.Port, cfg.HTTP.Host, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while starting http server: %s", err.Error())
	}
}
