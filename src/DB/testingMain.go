package main

import (
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	test "WhisperWave-BackEnd/src/DB/testing"
	"fmt"
)

func main() {
	// load AWS credentials config
	db_client := actionspkg.LoadDefaultConfig()

	// Initialize all the tables in DDB
	actionspkg.InitializeTables(db_client)

	actionspkg.InitChatHistory(db_client, "ChatHistory")
	actionspkg.InitUserAndGroupActions(db_client, "UserAndGroupInfo")
	actionspkg.InitUserServerMap(db_client, "UserServerMap")

	// Test tables
	fmt.Println("\n1. Testing \"UserServerMap\" Table .......")
	test.TestServerMap(db_client, "UserServerMap")
	fmt.Println("\n2. Testing \"ChatHistory\" Table .......")
	test.TestChatHistory(db_client, "ChatHistory")
	fmt.Println("\n3. Testing \"UserAndGroupInfo\" Table .......")
	test.TestUserAndGroupInfo(db_client, "UserAndGroupInfo")
}
