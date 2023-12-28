package routers

import (
	"WhisperWave-BackEnd/handlers"
	"WhisperWave-BackEnd/middleware"
	"WhisperWave-BackEnd/models"
	"net/http"

	"github.com/gorilla/mux"
)

func UserRouters(router *mux.Router, srv *models.Server) {
	router.HandleFunc("/", handlers.DefaultHandler)
	// router.HandleFunc("/ws", handlers.SingleUserChatHandler)
	router.Handle("/ws", middleware.WithServer(http.HandlerFunc(handlers.SingleUserChatHandler), srv))
}