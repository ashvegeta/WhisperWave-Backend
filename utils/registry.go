package utils

import (
	"WhisperWave-BackEnd/models"
	"log"
)

// var serverRegistry map[string]*models.Server
var UserRegistry map[string]*models.Server

func InitRegistry() {
	// serverRegistry = make(map[string]*models.Server)
	UserRegistry = make(map[string]*models.Server)	
}

//users
func GetServerForUser(userId string) (*models.Server) {
	srvInfo, exists := UserRegistry[userId]

	if !exists {
		log.Printf("user %s does not exist in user registry", userId)
	} 

	return srvInfo
}

func SetServerForUser(userId string, srv *models.Server) {
	UserRegistry[userId] = srv
}

// servers
func RegisterServer() {

}

func DeRegisterServer() {

}

func GetServerMetaData() {

}