package routers

import (
	"WhisperWave-BackEnd/handlers"
	"WhisperWave-BackEnd/middleware"
	server "WhisperWave-BackEnd/server"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router, srv *server.Server) {
	router.Handle("/ws", middleware.WithServer(http.HandlerFunc(handlers.SingleUserChatHandler), srv))
}
