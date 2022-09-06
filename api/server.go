package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/sbbullet/to-do/db"
	"github.com/sbbullet/to-do/logger"
	"github.com/sbbullet/to-do/token"
	"github.com/sbbullet/to-do/util"
)

type Server struct {
	config     *util.Config
	store      *db.Store
	router     *mux.Router
	tokenMaker token.Maker
}

func NewServer() *Server {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	config := util.LoadConfig("app", "env", wd)
	dbInstance := db.NewDB(config)
	store := db.NewStore(dbInstance)

	pasetoMaker, err := token.NewPasetoMaker(config.SymmetricKey)
	if err != nil {
		panic(err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: pasetoMaker,
	}

	// Setup server router
	server.setupRouter()

	return server
}

func (server *Server) setupRouter() {
	r := mux.NewRouter()
	r.Use(LoggingMiddleware())

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		util.RespondWithOk(w, "Yup, it's working. Explore the API documentation")
	})

	apiRoutes := r.PathPrefix("/api/v1/").Subrouter().StrictSlash(true)
	apiRoutes.HandleFunc("/users", server.RegisterUser).Methods(http.MethodPost)
	apiRoutes.HandleFunc("/users/login", server.LoginUser).Methods(http.MethodPost)

	authRoutes := apiRoutes.PathPrefix("/users").Subrouter().StrictSlash(true)
	authRoutes.Use(AuthMiddleware(server.tokenMaker))
	authRoutes.HandleFunc("/me", server.GetCurrentUser).Methods(http.MethodGet)

	server.router = r
}

func (server *Server) Run() {
	serverAddress := fmt.Sprintf("%s:%s", server.config.ServerHost, server.config.ServerPort)
	logger.Info(fmt.Sprintf("Server starting at http://%s", serverAddress))

	log.Fatal(http.ListenAndServe(serverAddress, server.router))
}
