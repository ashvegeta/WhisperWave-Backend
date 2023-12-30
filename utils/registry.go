package utils

import (
	"errors"
	"fmt"
)

// var serverRegistry map[string]*models.Server
var UserRegistry map[string]any

func InitRegistry() {
	// serverRegistry = make(map[string]*models.Server)
	UserRegistry = make(map[string]any)	
}

//users
func GetServerForUser(userId string) (any, error) {
	srvInfo, exists := UserRegistry[userId]
	var err error = nil

	if !exists {
		msg := fmt.Sprintf("\nuser %s does not exist in user registry", userId)
		err = errors.New(msg)
		return nil, err
	} 

	return srvInfo, err
}

func SetServerForUser(userId string, srv any) {
	UserRegistry[userId] = srv
}

// servers
func RegisterServer(srv any) {

}

func DeRegisterServer(srv any) {

}