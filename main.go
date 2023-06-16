package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"os"
)

func main() {
	server := &Server{
		port: os.Getenv("PORT"),
		publicEndpoints: []Endpoint{
			{
				Path:    "/token",
				Handler: Token,
				Methode: http.MethodPost,
			},
			{
				Path:    "/register",
				Handler: Register,
				Methode: http.MethodPost,
			},
		},
		privateEndpoints: []Endpoint{
			{
				Path:    "/user",
				Handler: GetUser,
				Methode: http.MethodGet,
			},
		},
		middlewares: []mux.MiddlewareFunc{JwtValidation()},
		publicRoot:  "/",
		privateRoot: "/api",
	}
	server.RunServer()
}
