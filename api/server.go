package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/util"
)

type Server struct {
	config *util.Config
	store  *db.Store
	router *mux.Router
}

func NewServer() *Server {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := util.LoadConfig("app", "env", wd)
	dbInstance := db.NewDB(config)
	store := db.NewStore(dbInstance)

	server := &Server{
		config: config,
		store:  store,
	}

	// Setup server router
	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.RespondWithOk(w, "Yup, it's working. Explore the API documentation")
	})

	apiRouter := r.PathPrefix("/api/v1/").Subrouter().StrictSlash(true)
	apiRouter.HandleFunc("/users", server.RegisterUser).Methods(http.MethodPost)

	server.router = r
}

func (server *Server) Run() {
	serverAddress := fmt.Sprintf("%s:%s", server.config.ServerHost, server.config.ServerPort)
	logger.Info(fmt.Sprintf("Server starting at http://%s", serverAddress))

	log.Fatal(http.ListenAndServe(serverAddress, server.router))
}
