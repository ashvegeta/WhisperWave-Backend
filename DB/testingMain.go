package main

import (
	subpkg "WhisperWave-BackEnd/DB/subpackage"
	test "WhisperWave-BackEnd/DB/testing"
	"context"
	"fmt"
	"log"

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

	// Test tables
	fmt.Println("\n1. Testing \"UserServerMap\" Table .......")
	test.TestServerMap(db_client, "UserServerMap")
	fmt.Println("\n2. Testing \"ChatHistory\" Table .......")
	test.TestChatHistory(db_client, "ChatHistory")
	fmt.Println("\n3. Testing \"UserAndGroupInfo\" Table .......")
	test.TestUserAndGroupInfo(db_client, "UserAndGroupInfo")
}
