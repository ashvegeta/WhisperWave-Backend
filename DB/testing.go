package main

import (
	subpkg "WhisperWave-BackEnd/DB/subpackage"
	"WhisperWave-BackEnd/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// load AWS credentials config
	config, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Println(err)
		return
	}

	db_client := subpkg.GetDBClient(config)

	// Initialize all the tables in DDB
	// subpkg.InitializeTables(db_client)

	// add new chat
	subpkg.AddNewChat(db_client, "ChatHistory", models.ChatHistory{
		PK:      "user1",
		SK:      fmt.Sprintf("%s-%d", string("user4"), time.Now().UnixMicro()),
		MID:     "mID1",
		MType:   "text/plain",
		Content: "hello, user, how are you ?",
	})

	// Perform Queries
	chatHistory, err := subpkg.LoadChatHistory(db_client, "ChatHistory", models.LoadChatInput{
		PK: "user1",
		SK: "user3",
	})

	if err != nil {
		log.Println(err)
		return
	}

	for _, chat := range chatHistory {
		fmt.Println(chat)
	}
}
