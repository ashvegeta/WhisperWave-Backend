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
	subpkg.InitializeTables(db_client)

	// 1. add new chat
	subpkg.AddNewChat(db_client, "ChatHistory", models.ChatHistory{
		PK:      "user1",
		SK:      fmt.Sprintf("%s-%d", string("user4"), time.Now().UnixMicro()),
		MID:     "mID1",
		MType:   "text/plain",
		Content: "hello, user, how are you ?",
	})

	// 2. add new user
	subpkg.AddNewUser(db_client, "UserAndGroupInfo", models.User{
		UserId:      "uid5",
		UserName:    "user5",
		Password:    "pwd",
		FriendsList: []string{"user1", "user2", "user3"},
		GroupList:   []string{"gid1", "gid2"},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// 3. add new group
	err = subpkg.AddNewGroup(db_client, "UserAndGroupInfo", models.Group{
		GroupId:   "gid2",
		GroupName: "group2",
		UserList:  []string{"user1, user2, user3, user4"},
	})
	if err != nil {
		log.Println(err)
		return
	}

	// Perform Queries
	// 1. Chat history
	chatHistory, err := subpkg.LoadChatHistory(db_client, "ChatHistory", models.ChatParams{
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

	// 2. User Info
	userInfo, err := subpkg.GetUserInfo(db_client, "UserAndGroupInfo", models.UserOrGroupParams{
		PK: "uid3",
	})
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(userInfo)

}
