package models

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	amqp "github.com/rabbitmq/amqp091-go"
)

// server metadata
type Server struct {
	Name string
	Addr string
	Mu	sync.Mutex
	Upgrader websocket.Upgrader
	ConnPool map[string] *websocket.Conn
	ConnLimit int
	MQ *MessageQueue
}

type MessageQueue struct {
	MQHost string //message queue address
	MQPort string 
	Queue amqp.Queue
	MQChannel *amqp.Channel
}

type UserServerMap struct {
	UserID string `json:"user_id"`
	ServerID string `json:"server"`
}

// default server config values
func (srv *Server) SetDefaultOps() {

}

// ------------ Server VARS ---------------- //
func (srv *Server) SetupServer(srvConfig map[string]string) {
	srv.Name = srvConfig["Name"]
	srv.Addr =  srvConfig["Addr"]

	R_SIZE, err := strconv.Atoi(srvConfig["ReadBufferSize"])
	if err != nil {
		fmt.Println(err)
		return
	}

	W_SIZE, err := strconv.Atoi(srvConfig["WriteBufferSize"])
	if err != nil {
		fmt.Println(err)
		return
	}

	srv.Upgrader =  websocket.Upgrader{
		ReadBufferSize: R_SIZE,
		WriteBufferSize: W_SIZE,
	}

	srv.ConnPool =  make(map[string]*websocket.Conn)

	connLimit, err := strconv.Atoi(srvConfig["ConnLimit"]) 
	if err != nil {
		fmt.Println("Error in string to int conversion: ", err)
		return
	}
	srv.ConnLimit = connLimit

	// message queue
	srv.MQ = &MessageQueue{}
	srv.MQ.MQHost = srvConfig["MQHost"]
	srv.MQ.MQPort = srvConfig["MQPort"]
	
	// Server.SetupMessageQueue()
}
// -----------------------------------------//

//Server read loop - Contains logic for handling user communication
func (srv* Server) ReadLoop(conn *websocket.Conn, userId string, IDgenerator func(string) string) {
	for {
		// read message from sender
		var recvMessage Message 
		err := conn.ReadJSON(&recvMessage)
		
		// in case of error in connection string, remove it from the pool and return error
		if err != nil {
			srv.Mu.Lock()
			delete(srv.ConnPool, userId)
			log.Println("Connection Error:", err)
			srv.Mu.Unlock()
			return
		}

		// generate a unique message ID
		recvMessage.MessageId = IDgenerator(recvMessage.SenderId)

		//get the sender userId
		receiverId := recvMessage.ReceiverId

		// check if the sender already has a connection with the server
		// IMPORTANT: IN FUTURE MODIFY THE CODE, AS 2 END USERS MAY NOT CONNECT TO SAME CHAT SERVER
		// Instead of connecting to the receiver, push it to message queue
		recvConn, exists := srv.ConnPool[receiverId]

		//if connection exists, send it to intended receiver
		if exists {
			if err := recvConn.WriteJSON(recvMessage); err != nil {
				log.Println("Server send Error:", err)
				return
			}
		}
	}
}

// setup the message queue for the server
func (srv* Server) SetupMessageQueue() {
	//host message queue as a service

	// connect with server
	addr := fmt.Sprintf("%s:%s",srv.MQ.MQHost, srv.MQ.MQPort)

	conn, err := amqp.Dial(addr); if err != nil {
		log.Println(err)
		return
	}

	ch, err := conn.Channel(); if err != nil {
		log.Println(err)
		return
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"MQName", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)

	if err!= nil {
		log.Println(err)
		return
	}

	srv.MQ.Queue = q
	srv.MQ.MQChannel = ch
}

// publish message to the target message queue (message queue of a different chat server)
func (srv* Server) PublishMessage(targetMQ MessageQueue, message []byte, mType string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err := targetMQ.MQChannel.PublishWithContext(ctx,
			"",     // exchange
			targetMQ.Queue.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType : mType,
				Body : message,
		})

	if err != nil {
		log.Println(err)
		return
	}
}

// consume message from queue
func (srv* Server) ConsumeMessage() (message []byte){
	var b []byte

	return b
}

// check heartbeat of connected websockets
func (srv* Server) HeartBeatMonitor() /*error*/ {
	/*for user, conn := range srv.ConnPool {
		conn.
	}*/
}