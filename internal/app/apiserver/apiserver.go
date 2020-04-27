package apiserver

import (
	"github.com/LKarlon/http-rest-api.git/internal/app/service"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

type APIServer struct {
	config            *Config
	logger            *logrus.Logger
	router            *mux.Router
	service  		  service.Service
}

func New(config *Config) *APIServer {
	return &APIServer{
		config:            config,
		logger:            logrus.New(),
		router:            mux.NewRouter(),
		service:           service.NewService(config.Store),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	s.configRouter()
	s.logger.Info("starting api server")
	go s.service.Worker()
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}

func (s *APIServer) configRouter() {
	s.router.HandleFunc("/inn", s.INN).Methods("GET")
}
func (s *APIServer) INN(w http.ResponseWriter, r *http.Request){
	s.service.GetInn(w, r)
}

