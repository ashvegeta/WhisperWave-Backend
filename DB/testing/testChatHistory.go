package testing

import (
	subpkg "WhisperWave-BackEnd/DB/subpackage"
	"WhisperWave-BackEnd/models"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestChatHistory(db_client *dynamodb.Client, tableName string) {
	// Init
	subpkg.InitChatHistory(db_client, tableName)

	// 1. add new chat
	uid1 := "uid1"
	uid2 := fmt.Sprintf("%s-%d", string("uid2"), time.Now().UnixMicro())
	err := subpkg.AddNewChat(models.ChatHistory{
		PK:      uid1,
		SK:      uid2,
		MID:     "mID1",
		MType:   "text/plain",
		Content: "hello, user, how are you ?",
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully inserted chat for user: %s\n", uid1)
	}

	// 2. load chat history
	chatHistory, err := subpkg.LoadChatHistory(models.ChatParams{
		PK: uid1,
		SK: uid2,
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully loaded chat of user: %s and user: %s\n", uid1, uid2)
	}

	for _, chat := range chatHistory {
		fmt.Println(chat)
	}

	// 3. Update an existing chat text
	updatedHistory, err := subpkg.UpdateChat(models.ChatParams{PK: uid1, SK: uid2}, models.ChatHistory{
		Content: "hello, user, how are you ? (modified)",
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully updated chat history", updatedHistory)
	}

	// 4. Delete a particular chat
	err = subpkg.DeleteChat(models.ChatParams{
		PK: uid1,
		SK: uid2,
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully deleted chat b/w user: %s and user: %s\n", uid1, uid2)
	}
}
