package app

import (
	"github.com/scmbr/test-task/internal/config"
	delivery "github.com/scmbr/test-task/internal/delivery/http"
	"github.com/scmbr/test-task/internal/server"
	"github.com/sirupsen/logrus"
)

func Run(configsDir string) {
	cfg, err := config.Init(configsDir)
	if err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
		return
	}
	srv := new(server.Server)
	handlers := delivery.NewHandler()
	if err := srv.Run(cfg.HTTP.Port, cfg.HTTP.Host, handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while starting http server: %s", err.Error())
	}
}
