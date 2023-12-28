package handlers

import (
	"net/http"

	"golang.org/x/net/websocket"
)

var groupConnections map[string] *websocket.Conn

func SubscribeToGroupHandler() {
	
}

func UnsubscribeFromGroupHandler() {
	
}

func GroupChatHandler(w http.ResponseWriter, r *http.Request) {
	// conn, err := upgrader.Upgrade(w, r, nil)

	// if err != nil {
	// 	log.Println("Error Upgrading Websocket: ", err)
	// 	return
	// }
	
	// //safely deffering the connection termination
	// defer conn.Close()

	// //fetch the user id from headers
	// userId := r.Header.Get("X-User-ID")
	// if userId == "" {
	// 	log.Println("User ID not found in headers. Closing connection.")
	// 	return
	// }

	// //map the user to respective connection and make sure it is sync locked
	// mu.Lock()
	// userConnections[userId] = conn
	// mu.Unlock()

	// for {
	// 	// read message from sender
	// 	var recvMessage models.Message 
	// 	err := conn.ReadJSON(&recvMessage)

	// 	// in case of error in connection string, remove it from the pool and return error
	// 	if err != nil {
	// 		mu.Lock()
	// 		delete(userConnections, userId)
	// 		log.Println("Connection Error:", err)
	// 		mu.Unlock()
	// 		return
	// 	}

	// 	fmt.Printf("server received message from %s : %s\n", recvMessage.SenderId, recvMessage.Content)

	// 	//get the sender userId
	// 	receiverId := recvMessage.ReceiverId

	// 	// check if the sender already has a connection with the server
	// 	// IMPORTANT: IN FUTURE MODIFY THE CODE, AS 2 END USERS MAY NOT CONNECT TO SAME CHAT SERVER
	// 	// Instead of connecting to the receiver, push it to message queue
	// 	recvConn, exists := userConnections[receiverId]

	// 	//if connection exists, send it to intended receiver
	// 	if exists {
	// 		if err := recvConn.WriteJSON(recvMessage); err != nil {
	// 			log.Println("Server send Error:", err)
	// 			return
	// 		}
	// 	}
	// }
}
