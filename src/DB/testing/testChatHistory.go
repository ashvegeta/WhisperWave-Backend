package testing

import (
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/models"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestChatHistory(db_client *dynamodb.Client, tableName string) {
	// Init
	actionspkg.InitChatHistory(db_client, tableName)

	// 1. add new user chat
	uid1 := "uid1"
	uid2 := fmt.Sprintf("%s-%d", string("uid2"), time.Now().UnixMicro())
	gid1 := "gid1"

	err := actionspkg.AddNewChat(models.ChatHistory{
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

	// add new group chat
	err = actionspkg.AddNewChat(models.ChatHistory{
		PK:      gid1,
		SK:      uid2,
		MID:     "mID1",
		MType:   "text/plain",
		Content: "hello, user, how are you ?",
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully inserted chat for group: %s\n", gid1)
	}

	// 2. load chat history
	chatHistory, err := actionspkg.LoadChatHistory(models.ChatParams{
		PK: uid1,
		// SK: uid2,
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
	updatedHistory, err := actionspkg.UpdateChat(models.ChatParams{PK: uid1, SK: uid2}, models.ChatHistory{
		Content: "hello, user, how are you ? (modified)",
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Successfully updated chat history", updatedHistory)
	}

	// 4. Delete a particular chat
	err = actionspkg.DeleteSingleChat(models.ChatParams{
		PK: uid1,
		SK: uid2,
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully deleted chat b/w user: %s and user: %s\n", uid1, uid2)
	}

	// // 5. delete user's group chat
	// err = actionspkg.DeleteUserGroupChat(models.ChatParams{
	// 	PK: "uid2",
	// })
	// if err != nil {
	// 	log.Println(err)

	// } else {
	// 	log.Printf("Successfully deleted user: %s,  group %s chat\n", uid2, gid1)
	// }
}
