package models

import (
	// "WhisperWave-BackEnd/utils"
	"WhisperWave-BackEnd/models"
	registry "WhisperWave-BackEnd/serviceRegistry"
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
	Name      string                     `json:"name"`
	Addr      string                     `json:"addr"`
	Mu        sync.Mutex                 `json:"mu"`
	Upgrader  websocket.Upgrader         `json:"upgrader"`
	ConnPool  map[string]*websocket.Conn `json:"conn_pool"`
	ConnLimit int                        `json:"conn_limit"`
	MQ        models.MessageQueue        `json:"mq"`
}

// default server config values
func (srv *Server) SetDefaultOps() {

}

// ------------ Server VARS ---------------- //
func (srv *Server) SetupServer(srvConfig map[string]any) {
	// configure server parameters
	srv.Name = srvConfig["Name"].(string)
	srv.Addr = srvConfig["Addr"].(string)
	srv.Upgrader = websocket.Upgrader{
		ReadBufferSize:  int(srvConfig["ReadBufferSize"].(float64)),
		WriteBufferSize: int(srvConfig["WriteBufferSize"].(float64)),
	}
	srv.ConnPool = make(map[string]*websocket.Conn)
	srv.ConnLimit = int(srvConfig["ConnLimit"].(float64))

	// configure server's message queue parameters
	srv.MQ = models.MessageQueue{}
	srv.MQ.MQName = srvConfig["MQName"].(string)
	srv.MQ.MQURI = srvConfig["MQURI"].(string)
	srv.MQ.MQParams = srvConfig["MQParams"].([]any)

	// create a queue OR get an active queue
	srv.GetQueue(srv.MQ)

	// start consuming as soon as the server starts
	go srv.ConsumeMessages()
}

// -----------------------------------------//

// Server read loop - Contains logic for handling user communication
func (srv *Server) ReadLoop(conn *websocket.Conn, userId string, IDgenerator func(string) string) {
	for {
		// read message from sender
		var recvMessage models.Message
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

		for _, receiverId := range recvMessage.ReceiverIds {
			// check if the receiver already has a connection with the server
			recvConn, exists := srv.ConnPool[receiverId]

			//if connection exists (i.e receiver connected to the same server), send it to intended receiver
			if exists {
				if err := recvConn.WriteJSON(recvMessage); err != nil {
					log.Println("Server send Error:", err)
					return
				}
			} else {
				// Send Message to the users who are connected to other chat servers
				ts, err := registry.GetServerForUser(receiverId)
				if err != nil {
					log.Println(err)
					return
				}

				copyMessage := recvMessage
				copyMessage.ReceiverIds = []string{receiverId}
				targetSrv := ts.(models.UserServerMap)

				var targetMQ models.MessageQueue = targetSrv.ServerInfo.MQ
				srv.PublishMessage(targetMQ, copyMessage, copyMessage.MessageType)
			}
		}
	}
}

// setup the message queue for the server
func (srv *Server) GetQueue(targetMQ models.MessageQueue) *amqp.Channel {
	// connect with server
	conn, err := amqp.Dial(targetMQ.MQURI)
	if err != nil {
		log.Println(err)
		return nil
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Println(err)
		return nil
	}

	// if the queue already exists, with below params its fetches the same
	_, err = ch.QueueDeclare(
		targetMQ.MQName,             // name
		targetMQ.MQParams[0].(bool), // durable
		targetMQ.MQParams[1].(bool), // delete when unused
		targetMQ.MQParams[2].(bool), // exclusive
		targetMQ.MQParams[3].(bool), // no-wait
		nil,                         // arguments
	)

	if err != nil {
		log.Println(err)
		return nil
	}

	return ch
}

// publish message to the target message queue (message queue of a different chat server)
func (srv *Server) PublishMessage(targetMQ models.MessageQueue, message models.Message, mType string) {
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
		"",              // exchange
		targetMQ.MQName, // routing key
		false,           // mandatory
		false,           // immediate
		amqp.Publishing{
			ContentType: mType,
			Body:        mBody,
		})

	if err != nil {
		log.Println("Failed to publish message :", err)
		return
	}
}

// send message to intended user
func (srv *Server) SendMessage(message models.Message) {
	for _, receiverId := range message.ReceiverIds {
		// check if user is currently online/connected
		conn, exists := srv.ConnPool[receiverId]

		// if the message recepient is connected to the current server, then forward the message
		if exists {
			if err := conn.WriteJSON(message); err != nil {
				log.Println("Error in forwaring message to User : ", receiverId, err)
				return
			}
		} else /*handle disconnected users properly*/ {
			log.Printf("\nUser %s does not have an active connection", receiverId)
			return
		}
	}
}

// consume message from queue
func (srv *Server) ConsumeMessages() {
	ch := srv.GetQueue(srv.MQ)
	var intermData models.Message

	msgs, err := ch.Consume(
		srv.MQ.MQName, // queue
		srv.Addr,      // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		nil,           // args
	)

	if err != nil {
		panic(err)
	}

	for {
		// try to re-setablish connection with the MQ to start consuming messages (IN CASE ITS CLOSED)
		if ch.IsClosed() || ch == nil {
			ch = srv.GetQueue(srv.MQ)

			msgs, err = ch.Consume(
				srv.MQ.MQName, // queue
				srv.Addr,      // consumer
				true,          // auto-ack
				false,         // exclusive
				false,         // no-local
				false,         // no-wait
				nil,           // args
			)

			if err != nil {
				panic(err)
			}
		}

		for msg := range msgs {
			json.Unmarshal(msg.Body, &intermData)
			fmt.Printf("\n inside [%s] : %s", srv.Name, intermData)
			srv.SendMessage(intermData)
		}
	}
}

// check heartbeat of connected websockets
func (srv *Server) HeartBeatMonitor() /*error*/ {
	/*for user, conn := range srv.ConnPool {
		conn.
	}*/
}
