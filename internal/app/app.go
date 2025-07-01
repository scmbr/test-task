package app

import (
	delivery "github.com/scmbr/test-task/internal/delivery/http"
	"github.com/scmbr/test-task/internal/server"
	"github.com/sirupsen/logrus"
)

func Run(configsDir string) {
	srv := new(server.Server)
	handlers := delivery.NewHandler()
	if err := srv.Run("8080", handlers.InitRoutes()); err != nil {
		logrus.Fatalf("error occurred while starting http server: %s", err.Error())
	}
}
