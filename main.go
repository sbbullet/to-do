package main

import (
	"github.com/sbbullet/to-do/api"
	"github.com/sbbullet/to-do/logger"
)

func main() {

	// Initialize logger
	logger.InitializeLogger()

	// Create a new server
	server := api.NewServer()

	// Run the server
	server.Run()
}
