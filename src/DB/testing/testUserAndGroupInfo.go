package testing

import (
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestUserAndGroupInfo(db_client *dynamodb.Client, tableName string) {
	// Init
	actionspkg.InitUserAndGroupActions(db_client, tableName)

	// Test
	uid1 := "uid1"
	gid1 := "gid1"

	// 1. add new user
	err := actionspkg.AddNewUserOrGroup(models.User{
		UserId:      uid1,
		UserName:    "user1",
		Password:    "pwd",
		EmailID:     "user1@gmail.com",
		FriendsList: []string{"uid2", "uid3", "uid4"},
		GroupList:   []string{"gid1", "gid2"},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully inserted new user: %s\n", uid1)
	}

	err = actionspkg.AddNewUserOrGroup(models.User{
		UserId:      "uid2",
		UserName:    "user2",
		Password:    "pwd",
		EmailID:     "user2@gmail.com",
		FriendsList: []string{"uid1", "uid3", "uid4"},
		GroupList:   []string{"gid1", "gid2"},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully inserted new user: %s\n", "uid2")
	}

	// 2. add new group
	err = actionspkg.AddNewUserOrGroup(models.Group{
		GroupId:   gid1,
		GroupName: "group1",
		UserList:  []string{"uid1", "uid2", "uid3", "uid4"},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully inserted new group: %s\n", gid1)
	}

	// 3. get User Info
	userInfo, err := actionspkg.GetUserInfo(models.UserOrGroupParams{
		PK: uid1,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully fetched user info: %s\n", gid1)
		fmt.Println(userInfo)
	}

	// 4. get Group Info
	groupInfo, err := actionspkg.GetGroupInfo(models.UserOrGroupParams{
		PK: gid1,
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully fetched group info: %s\n", gid1)
		fmt.Println(groupInfo)
	}

	// 3. Update User Info
	updatedInfo, err := actionspkg.UpdateUserOrGroupInfo(models.UserOrGroupParams{PK: uid1}, models.User{
		Password:    "pwd2",
		FriendsList: []string{"uid2", "uid3", "uid5"},
	})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully updated user info: %s\n", uid1)
		fmt.Println(updatedInfo)
	}

	// 4. Delete User/Group Info
	deletedInfo, err := actionspkg.DeleteUserOrGroup(models.UserOrGroupParams{PK: uid1})
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Successfully deleted user info: %s\n", uid1)
		fmt.Println(deletedInfo)
	}
}
