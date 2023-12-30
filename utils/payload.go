package utils

// import (
// 	"WhisperWave-BackEnd/models"
// 	"encoding/json"
// 	"log"
// )

// func SerializeJSON(data models.Message) []byte {
// 	serData, err := json.Marshal(data)

// 	if err != nil {
// 		log.Println("Error serializing data:", serData)
// 		log.Println(err)
// 		return nil
// 	}

// 	return serData
// }

// func DeserializeJSON(encodedData []byte, pointerData *models.Message) {
// 	err := json.Unmarshal(encodedData, pointerData)

// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
//  }