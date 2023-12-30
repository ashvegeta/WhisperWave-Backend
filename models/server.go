package models

import (
	"WhisperWave-BackEnd/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	amqp "github.com/rabbitmq/amqp091-go"
)

// server metadata
type Server struct {
	Name string								`json:"name"`
	Addr string								`json:"addr"`
	Mu	sync.Mutex							`json:"mu"`
	Upgrader websocket.Upgrader				`json:"upgrader"`
	ConnPool map[string] *websocket.Conn	`json:"conn_pool"`
	ConnLimit int							`json:"conn_limit"`
	MQ MessageQueue							`json:"mq"`
}

type MessageQueue struct {
	MQName string	`json:"mq_name"`
	MQURI string 	`json:"mq_uri"`
	MQParams []any	`json:"mq_params"`
}

type UserServerMap struct {
	UserID string `json:"user_id"`
	ServerID string `json:"server"`
}

// default server config values
func (srv *Server) SetDefaultOps() {

}

// ------------ Server VARS ---------------- //
func (srv *Server) SetupServer(srvConfig map[string]any) {
	// configure server parameters
	srv.Name = srvConfig["Name"].(string)
	srv.Addr =  srvConfig["Addr"].(string)
	srv.Upgrader =  websocket.Upgrader{
		ReadBufferSize: int(srvConfig["ReadBufferSize"].(float64)),
		WriteBufferSize: int(srvConfig["WriteBufferSize"].(float64)),
	}
	srv.ConnPool =  make(map[string]*websocket.Conn)
	srv.ConnLimit = int(srvConfig["ConnLimit"].(float64))

	// configure server's message queue parameters
	srv.MQ = MessageQueue{}
	srv.MQ.MQName = srvConfig["MQName"].(string)
	srv.MQ.MQURI = srvConfig["MQURI"].(string)
	srv.MQ.MQParams = srvConfig["MQParams"].([]any)
	
	// create a queue OR get an active queue
	srv.GetQueue(srv.MQ)

	// start consuming as soon as the server starts
	// go srv.ReadMQLoop()
}
// -----------------------------------------//

//Server read loop - Contains logic for handling user communication
func (srv* Server) ReadLoop(conn *websocket.Conn, userId string, IDgenerator func(string) string) {
	for {
		fmt.Println("inside loop")
		// read message from sender
		var recvMessage Message 
		err := conn.ReadJSON(&recvMessage)
		
		// in case of error in reading message, remove it from the pool and return error
		if err != nil {
			srv.Mu.Lock()
			delete(srv.ConnPool, userId)
			log.Println("Connection Error:", err)
			srv.Mu.Unlock()
			return
		}

		// generate a unique message ID
		recvMessage.MessageId = IDgenerator(recvMessage.SenderId)

		//get the receiver's userId
		receiverId := recvMessage.ReceiverId

		// check if the sender already has a connection with the server
		// IMPORTANT: IN FUTURE MODIFY THE CODE, AS 2 END USERS MAY NOT CONNECT TO SAME CHAT SERVER
		// Instead of connecting to the receiver, push it to message queue
		recvConn, exists := srv.ConnPool[receiverId]

		//if connection exists (i.e receiver connected to the same server), send it to intended receiver
		if exists {
			if err := recvConn.WriteJSON(recvMessage); err != nil {
				log.Println("Server send Error:", err)
				return
			}
		} else {
			// fmt.Print(utils.GetServerForUser(receiverId))
			targetSrv, err := utils.GetServerForUser(receiverId)
			if err != nil {
				log.Println(err)
				return
			}

			var targetMQ MessageQueue = targetSrv.(*Server).MQ
			srv.PublishMessage(targetMQ, recvMessage, recvMessage.MessageType)	
		}
	}
}

// setup the message queue for the server
func (srv* Server) GetQueue(targetMQ MessageQueue) *amqp.Channel {
	//host message queue as a service

	// connect with server
	conn, err := amqp.Dial(targetMQ.MQURI); if err != nil {
		log.Println(err)
		return nil
	}

	ch, err := conn.Channel(); if err != nil {
		log.Println(err)
		return nil
	}

	// if the queue already exists, with below params its fetches the same
	_, err = ch.QueueDeclare(
		targetMQ.MQName, 			   // name
		targetMQ.MQParams[0].(bool),   // durable
		targetMQ.MQParams[1].(bool),   // delete when unused
		targetMQ.MQParams[2].(bool),   // exclusive
		targetMQ.MQParams[3].(bool),   // no-wait
		nil,     			  		   // arguments
	)

	if err!= nil {
		log.Println(err)
		return nil
	}

	return ch
}

// publish message to the target message queue (message queue of a different chat server)
func (srv* Server) PublishMessage(targetMQ MessageQueue, message Message, mType string) {
	// create context for publishing
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// get queue for publishing
	ch := srv.GetQueue(targetMQ)

	// encode message
	mBody, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}

	// publish message
	err = ch.PublishWithContext(ctx,
			"",     		 // exchange
			targetMQ.MQName, // routing key
			false,  		 // mandatory
			false,  		 // immediate
			amqp.Publishing{
				ContentType : mType,
				Body : mBody,
		})

	if err != nil {
		log.Println("Failed to publish message :", err)
		return
	}
}

// send message to intended user
func (srv *Server) SendMessage(message Message) {
	receiverId := message.ReceiverId

	// check if user is currently online/connected
	conn, exists := srv.ConnPool[receiverId]

	// if the message recepient is connected to the current server, then forward the message
	if exists {
		if err := conn.WriteJSON(message); err != nil {
			log.Println("Error in forwaring message to User : ", receiverId, err)
			return
		} 
	} /*handle disconnected users properly*/ else {
		log.Printf("\nUser %s does nat have an active connection", receiverId)
		return
	}
}

// consume message from queue
func (srv* Server) ConsumeMessages() {
	for	{	

		ch := srv.GetQueue(srv.MQ)
		
		msgs, err := ch.Consume(
			srv.MQ.MQName, // queue
			"",        // consumer
			true,      // auto-ack
			false,     // exclusive
			false,     // no-local
			false,     // no-wait
			nil,       // args
		)

		// defer ch.Close()
		
		if err != nil {
			panic(err)
		}
		
		var intermData Message

		for msg := range msgs {
			json.Unmarshal(msg.Body, &intermData)
			fmt.Printf("\n inside [%s] : %s", srv.Name, intermData)
			srv.SendMessage(intermData)
		}
	}
}

// check heartbeat of connected websockets
func (srv* Server) HeartBeatMonitor() /*error*/ {
	/*for user, conn := range srv.ConnPool {
		conn.
	}*/
}