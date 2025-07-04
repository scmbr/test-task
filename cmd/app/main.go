package main

import (
	_ "github.com/scmbr/test-task/docs"
	"github.com/scmbr/test-task/internal/app"
)

const configsDir = "configs"

// @title Authorization API test task
// @version 1.0
// @description Реализация тестового задания medods
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @host localhost:8080
// @BasePath /
func main() {
	app.Run(configsDir)
}
