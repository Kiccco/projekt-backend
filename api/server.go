package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type Server struct {
	*mux.Router
}

func NewServer() *Server {
	server := &Server{mux.NewRouter()}

	auth := server.PathPrefix("/account").Subrouter()
	auth.HandleFunc("/login", handleLogin)
	auth.HandleFunc("/register", handleRegister)

	go http.ListenAndServe(":8000", server)
	return server
}

func handleLogin(writer http.ResponseWriter, req *http.Request) {

}

func handleRegister(writer http.ResponseWriter, req *http.Request) {

}
