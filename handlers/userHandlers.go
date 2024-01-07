package handlers

import (
	"WhisperWave-BackEnd/models"
	"WhisperWave-BackEnd/utils"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

func RandomGenerator(seed int64) string {
	rand.Seed(seed)
	randomNumber := rand.Intn(1000000)
	return strconv.Itoa(randomNumber)
}

func DefaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello world"))
}

func SendFriendRequestHandler() {

}

func CancelFriendRequestHandler() {

}

// you can use channels
func SingleUserChatHandler(w http.ResponseWriter, r *http.Request) {
	// get server metadata from context
	Server, ok := r.Context().Value(models.ServerContext{Key: "server"}).(*models.Server)
	if !ok {
		http.Error(w, "Server not found in context", http.StatusInternalServerError)
		return
	}

	//fetch the user id from headers
	userId := r.Header.Get("X-User-ID")
	if userId == "" {
		log.Println("User ID not found in headers. Closing connection.")
		return
	}

	// check for connection limit
	if len(Server.ConnPool) == Server.ConnLimit {
		w.WriteHeader(http.StatusTooManyRequests)
		w.Write([]byte("Error : Max Connection Limit Exceeded"))
	}

	//upgrade the HTTP connnection to websocket (IF not done already)
	var conn *websocket.Conn
	conn, exists := Server.ConnPool[userId]

	if !exists {
		var err error
		conn, err = Server.Upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println("Error Upgrading Websocket: ", err)
			return
		}
		//safely defer the connection termination
		defer conn.Close()

		//map the user to respective connection and make sure it synchronous
		Server.Mu.Lock()
		Server.ConnPool[userId] = conn
		utils.SetServerForUser(userId, Server)
		fmt.Println(utils.UserRegistry)
		Server.Mu.Unlock()
	} else {
		log.Printf("\nUser %s is already connected to the chat server", userId)
	}

	// read loop
	Server.ReadLoop(conn, userId, utils.GenerateMessageID)
}
