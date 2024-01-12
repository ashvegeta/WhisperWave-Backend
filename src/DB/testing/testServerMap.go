package testing

import (
	subpkg "WhisperWave-BackEnd/DB/actionspkg"
	"WhisperWave-BackEnd/models"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestServerMap(db_client *dynamodb.Client, tableName string) {
	// Init
	subpkg.InitServerMap(db_client, tableName)

	// Call Actions
	uid := "uid1"
	// put
	err := subpkg.PutServerMap(models.UserServerMap{
		UserID: uid,
		ServerInfo: models.ServerInfo{
			SrvName: "server1",
			SrvAddr: "localhost:8080",
			MQ: models.MessageQueue{
				MQName:   "MQ1",
				MQURI:    "localhost:5000",
				MQParams: []any{true, false, false, false},
			},
		},
	})
	if err != nil {
		log.Println(err)

	} else {
		log.Printf("Successfully inserted user : %s\n", uid)
	}

	// // get
	items, err := subpkg.GetServerMap(uid)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Successfully fetched users")
		log.Println(items)
	}

	// delete
	err = subpkg.DeleteServerMap(uid)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("Successfully deleted user : %s\n", uid)
	}
}
