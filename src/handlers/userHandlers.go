package handlers

import (
	server "WhisperWave-BackEnd/server"
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/models"
	registry "WhisperWave-BackEnd/src/serviceRegistry"
	"WhisperWave-BackEnd/src/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"

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

// load user info
func GetUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	userId := r.Header.Get("X-User-ID")
	if userId == "" {
		log.Println("User ID not found in headers. Closing connection.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	userInfo, err := actionspkg.GetUserInfo(models.UserOrGroupParams{PK: userId})
	if err != nil {
		log.Println("error fetching user info : ", err)
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}

	json.NewEncoder(w).Encode(userInfo)
}

// load  chat for a user/group
func LoadChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	// ID1 - userid or groupid (PK in ddb table)
	// ID2 - the SK in the ddb table (i.e the userid of the person initiating conversation in group, or the recepient in a private chat)
	Id1 := r.Header.Get("ID1")
	Id2 := r.Header.Get("ID2")

	if Id1 == "" {
		log.Println("intended ID not found in header. Closing connection.")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// fetch chat history for private (uid, uid-ts) chat
	var messages []models.Message
	chatParams := models.ChatParams{PK: Id1}
	if Id2 != "" {
		chatParams.SK = Id2
	}

	chatHistory, err := actionspkg.LoadChatHistory(chatParams)
	if err != nil {
		log.Println("error loading chat history : ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	for _, chat := range chatHistory {
		SK := strings.Split(chat.SK, "-")

		messages = append(messages, models.Message{
			SenderId:    chat.PK,
			MessageId:   chat.MID,
			ReceiverIds: []string{SK[0]},
			MessageType: chat.MType,
			Content:     chat.Content,
			TimeStamp:   SK[1],
		})
	}

	json.NewEncoder(w).Encode(messages)
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

	// read loop
	Server.ReadLoop(conn, userId, utils.GenerateMessageID)
}
