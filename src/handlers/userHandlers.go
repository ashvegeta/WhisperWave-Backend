package handlers

import (
	server "WhisperWave-BackEnd/server"
	"WhisperWave-BackEnd/src/models"
	registry "WhisperWave-BackEnd/src/serviceRegistry"
	"WhisperWave-BackEnd/src/utils"
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

// load  chat user
func LoadChatHistoryHandler(w http.ResponseWriter, r *http.Request) {

}

// you can use channels
func SingleUserChatHandler(w http.ResponseWriter, r *http.Request) {
	// get server metadata from context
	Server, ok := r.Context().Value(models.ServerContext{Key: "server"}).(*server.Server)
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
		registry.SetServerForUser(userId, models.ServerInfo{
			SrvName: Server.Name,
			SrvAddr: Server.Addr,
			MQ:      Server.MQ,
		})
		fmt.Println(registry.GetServerForUser(userId))
		Server.Mu.Unlock()
	} else {
		log.Printf("\nUser %s is already connected to the chat server", userId)
	}

	// load chat history
	chatHistory := Server.LoadChatHistory(userId)

	for _, chat := range chatHistory {
		conn.WriteJSON(chat)
	}

	// read loop
	Server.ReadLoop(conn, userId, utils.GenerateMessageID)
}
