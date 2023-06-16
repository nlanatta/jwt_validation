package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Endpoint struct {
	Path    string
	Handler http.HandlerFunc
	Methode string
}

type Server struct {
	port             string
	publicEndpoints  []Endpoint
	privateEndpoints []Endpoint
	publicRoot       string
	privateRoot      string
	middlewares      []mux.MiddlewareFunc
}

func (s *Server) RunServer() {
	r := mux.NewRouter()
	basePath := r.PathPrefix(s.publicRoot).Subrouter()
	api := basePath.PathPrefix(s.privateRoot).Subrouter()

	s.registerPublicEndpoints(basePath)
	s.registerEndpoints(api)

	for _, v := range s.middlewares {
		api.Use(v)
	}

	log.Fatal(http.ListenAndServe(":"+s.port, r))
}

func (s *Server) registerEndpoints(api *mux.Router) {
	for _, v := range s.privateEndpoints {
		api.Handle(v.Path, v.Handler).Methods(v.Methode)
	}
}

func (s *Server) registerPublicEndpoints(api *mux.Router) {
	for _, v := range s.publicEndpoints {
		api.Handle(v.Path, v.Handler).Methods(v.Methode)
	}
}
