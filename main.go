package main

import (
	"github.com/sbbullet/to-do/api"
	"github.com/sbbullet/to-do/logger"
)

func main() {
	logger.InitializeLogger()
	logger.Info("Starting app")

	server := api.NewServer()
	server.Run()
}
