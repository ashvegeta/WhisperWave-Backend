package routers

import (
	server "WhisperWave-BackEnd/server"
	"WhisperWave-BackEnd/src/handlers"

	"github.com/gorilla/mux"
)

func InitRouter(router *mux.Router, srv *server.Server) {
	// init common routers
	router.HandleFunc("/check", handlers.TokenHandler).Methods("POST")
	router.HandleFunc("/", handlers.DefaultHandler).Methods("GET")
	router.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	router.HandleFunc("/signup", handlers.SignupHandler).Methods("POST")

	// init user routers
	UserRouter(router, srv)
}
