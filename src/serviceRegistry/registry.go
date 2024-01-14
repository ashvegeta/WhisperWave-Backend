package serviceregistry

import (
	subpkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/models"
	"errors"
	"fmt"
	"log"
)

// users
func GetServerForUser(userId string) (any, error) {
	srvInfo, err := subpkg.GetServerMap(userId)

	if len(srvInfo) == 0 {
		msg := fmt.Sprintf("\nuser %s does not exist in user registry", userId)
		err = errors.New(msg)
		return nil, err
	}

	return srvInfo[0], err
}

func SetServerForUser(userId string, srv any) {
	err := subpkg.PutServerMap(models.UserServerMap{
		UserID:     userId,
		ServerInfo: srv.(models.ServerInfo),
	})

	if err != nil {
		log.Println(err)
		return
	}
}

// servers
func RegisterServer(srv any) {

}

func DeRegisterServer(srv any) {

}
