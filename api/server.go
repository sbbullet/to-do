package api

import (
	"log"
	"os"

	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/util"
)

type Server struct {
	config *util.Config
	store  *db.Store
}

func NewServer() *Server {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := util.LoadConfig("app", "env", wd)
	log.Print(config)
	// store := db.NewStore()

	server := &Server{
		config: config,
		// store:  store,
	}

	return server
}

func (server *Server) Run() {
	log.Printf("Hello from runner")
}
