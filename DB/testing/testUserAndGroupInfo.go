package testing

import (
	subpkg "WhisperWave-BackEnd/DB/subpackage"
	"WhisperWave-BackEnd/models"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestUserAndGroupInfo(db_client *dynamodb.Client, tableName string) {
	// Init
	subpkg.InitUserAndGroupActions(db_client, tableName)

	// Test
	uid1 := "uid1"
	gid1 := "gid1"

	// 1. add new user
	err := subpkg.AddNewUserOrGroup(models.User{
		UserId:      uid1,
		UserName:    "user1",
		Password:    "pwd",
		FriendsList: []string{"user2", "user3", "user4"},
		GroupList:   []string{"gid1", "gid2"},
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully inserted new user: %s\n", uid1)
	}

	// 2. add new group
	err = subpkg.AddNewUserOrGroup(models.Group{
		GroupId:   gid1,
		GroupName: "group2",
		UserList:  []string{"user1, user2, user3, user4"},
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully inserted new group: %s\n", gid1)
	}

	// 3. get User Info
	userInfo, err := subpkg.GetUserInfo(models.UserOrGroupParams{
		PK: uid1,
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully fetched user info: %s\n", gid1)
		fmt.Println(userInfo)
	}

	// 4. get Group Info
	groupInfo, err := subpkg.GetGroupInfo(models.UserOrGroupParams{
		PK: gid1,
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully fetched group info: %s\n", gid1)
		fmt.Println(groupInfo)
	}

	// 3. Update User Info
	updatedInfo, err := subpkg.UpdateUserOrGroupInfo(models.UserOrGroupParams{PK: uid1}, models.User{
		Password:    "pwd2",
		FriendsList: []string{"user1", "user2", "user5"},
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully updated user info: %s\n", uid1)
		fmt.Println(updatedInfo)
	}

	// 4. Delete User/Group Info
	deletedInfo, err := subpkg.DeleteUserOrGroup(models.UserOrGroupParams{PK: gid1})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully deleted group info: %s\n", gid1)
		fmt.Println(deletedInfo)
	}
}
