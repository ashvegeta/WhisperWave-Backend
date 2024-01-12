package main

import (
	subpkg "WhisperWave-BackEnd/DB/actionspkg"
	test "WhisperWave-BackEnd/DB/testing"
	"fmt"
)

func main() {
	// load AWS credentials config
	db_client := subpkg.LoadDefaultConfig()

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
