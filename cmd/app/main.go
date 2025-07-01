package main

import "github.com/scmbr/test-task/internal/app"

const configsDir = "configs"

func main() {
	app.Run(configsDir)
}
