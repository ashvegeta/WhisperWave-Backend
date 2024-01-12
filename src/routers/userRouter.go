package routers

import (
	server "WhisperWave-BackEnd/server"
	"WhisperWave-BackEnd/src/handlers"
	"WhisperWave-BackEnd/src/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router, srv *server.Server) {
	router.Handle("/ws", middleware.WithServer(http.HandlerFunc(handlers.SingleUserChatHandler), srv))
}
