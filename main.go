package main

import (
	"WhisperWave-BackEnd/models"
	"WhisperWave-BackEnd/routers"
	"WhisperWave-BackEnd/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func InitServer(srvInfo map[string]interface{}) *models.Server {
	var srv *models.Server = &models.Server{}
	
	var configMap map[string]string = map[string]string{
		"Name" : srvInfo["SRV_NAME"].(string),
		"Addr" : srvInfo["SRV_HOST"].(string) + ":" +srvInfo["SRV_PORT"].(string),
		"ReadBufferSize" : srvInfo["SRV_BUFFER_READ_SIZE"].(string),
		"WriteBufferSize" : srvInfo["SRV_BUFFER_WRITE_SIZE"].(string),
		"ConnLimit" : srvInfo["SRV_CONN_LIMIT"].(string),
		"MQHost" : srvInfo["MQ_HOST"].(string),
		"MQPort" : srvInfo["MQ_PORT"].(string),
	}

	srv.SetupServer(configMap)
	// fmt.Println(configMap)

	return srv
}

func StartServer(chatServer *models.Server) {
	// setup routers
	r := mux.NewRouter()
	routers.UserRouters(r, chatServer)

	// setup http server
	httpsrv := &http.Server{
		Handler: r,
		Addr: chatServer.Addr,
	}

	fmt.Printf("listening at address: %s\n", httpsrv.Addr)
	httpsrv.ListenAndServe()
}

func main() {
	// initialize variables
	var (
		serversInfo []interface{} = utils.LoadSrvConfig()
		wg sync.WaitGroup
	)
	
	// setup registry
	utils.InitRegistry()

	// initialize servers
	for c, srvInfo := range serversInfo {
		wg.Add(c + 1)
		
		go func(srvInfo interface{}) {
			server := InitServer(srvInfo.(map[string]interface{}))
			StartServer(server)
			wg.Done()
		}(srvInfo)
	}

	// wait for all the servers to initialize before ending 
	wg.Wait()
}