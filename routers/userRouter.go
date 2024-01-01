package routers

import (
	"WhisperWave-BackEnd/handlers"
	"WhisperWave-BackEnd/middleware"
	"WhisperWave-BackEnd/models"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRouter(router *mux.Router, srv *models.Server) {
	router.Handle("/ws", middleware.WithServer(http.HandlerFunc(handlers.SingleUserChatHandler), srv))
}