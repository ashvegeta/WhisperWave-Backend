package testing

import (
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/models"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func TestServerMap(db_client *dynamodb.Client, tableName string) {
	// Init
	actionspkg.InitUserServerMap(db_client, tableName)

	// Call Actions
	uid := "uid1"
	// put
	err := actionspkg.PutServerMap(models.UserServerMap{
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

	// get
	items, err := actionspkg.GetServerMap(uid)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Println("Successfully fetched users")
		log.Println(items)
	}

	// delete
	err = actionspkg.DeleteServerMap(uid)
	if err != nil {
		log.Println(err)
		return
	} else {
		log.Printf("Successfully deleted user : %s\n", uid)
	}
}
