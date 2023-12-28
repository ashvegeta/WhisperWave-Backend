package main

import (
	"WhisperWave-BackEnd/models"
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

//  ----------- CLIENT VARS ---------- //
var(
	mu sync.Mutex
)

func RandomGenerator(seed int64) string {
	rand.Seed(seed)
	randomNumber := rand.Intn(1000000)
	return strconv.Itoa(randomNumber)
}

func TestClient() {
	scan := bufio.NewScanner(os.Stdin)

	var seed int64

	fmt.Print("Enter your unique seed: ")
	if scan.Scan() {
		value, err := strconv.Atoi(scan.Text())
		if err != nil {
			fmt.Println(err)
		}

		seed = int64(value)
	}
	
	// Set Headers
	header := http.Header{}
	var receiverId, senderId string
	
	if seed == 1 {
		senderId = RandomGenerator(1)
		receiverId = RandomGenerator(2)
	} else {
		senderId = RandomGenerator(2)
		receiverId = RandomGenerator(1)
	}
	
	header.Set("X-User-ID", senderId)
	
	// set URL
	var u url.URL

	if seed == 1 {
		u = url.URL{
			Scheme : "ws",
			Host: "localhost:8080",
			Path: "/ws",
		}
	} else {
		u = url.URL{
			Scheme : "ws",
			Host: "localhost:8081",
			Path: "/ws",
		}
	}

	// dial for a websocket connection
	dialer := websocket.Dialer{
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

	conn, resp , err := dialer.Dial(u.String(), header)

	if err != nil {
		log.Println("Error Dialing to the websocket", u.Host, " : " ,err)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		} 

		// Print the response body as a string
		fmt.Println(string(body))
		return
	}

	//Read Loop
	for {
		// input message
		var sentMessage, recvMessage models.Message
		var txtMsg string
		
		//receive message (in a GO routine)
		go func() {
			err2 := conn.ReadJSON(&recvMessage)
			if err2 != nil {
				log.Println("Error in receiving message : ", err)
				return
			}
	
			//print received message
			fmt.Printf("\n%s : %s\n", recvMessage.SenderId , recvMessage.Content)
		}()

		// send message
		fmt.Printf("You [%s]: ", senderId)
		if scan.Scan() {
			txtMsg = scan.Text()
		}

		sentMessage = models.Message{
			SenderId: senderId,
			ReceiverId: receiverId,
			Content: txtMsg,
			MessageType: "text",
			TimeStamp: time.Now(),
		}
	
		mu.Lock()
		err1 := conn.WriteJSON(sentMessage)
		mu.Unlock()

		if err1 != nil {
			log.Println("Error in sending message : ", err)
			return
		}
	}
}

func main() {
	TestClient()
}