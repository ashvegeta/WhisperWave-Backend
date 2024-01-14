package main

import (
	server "WhisperWave-BackEnd/server"
	actionspkg "WhisperWave-BackEnd/src/DB/actionspkg"
	"WhisperWave-BackEnd/src/routers"
	"WhisperWave-BackEnd/src/utils"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

func InitServer(srvInfo map[string]interface{}) *server.Server {
	var srv *server.Server = &server.Server{}

	var configMap map[string]any = map[string]any{
		"Name":            srvInfo["SRV_NAME"],
		"Addr":            srvInfo["SRV_HOST"].(string) + ":" + srvInfo["SRV_PORT"].(string),
		"ReadBufferSize":  srvInfo["SRV_BUFFER_READ_SIZE"],
		"WriteBufferSize": srvInfo["SRV_BUFFER_WRITE_SIZE"],
		"ConnLimit":       srvInfo["SRV_CONN_LIMIT"],
		"MQName":          srvInfo["MQ_NAME"],
		"MQURI":           srvInfo["MQ_URI"],
		"MQParams":        srvInfo["MQ_PARAMS"],
	}

	srv.SetupServer(configMap)

	return srv
}

func StartServer(chatServer *server.Server) {
	// setup routers
	r := mux.NewRouter()
	routers.InitRouter(r, chatServer)

	// setup http server
	httpsrv := &http.Server{
		Handler: r,
		Addr:    chatServer.Addr,
	}

	// start server
	fmt.Printf("listening at address: %s\n", httpsrv.Addr)
	httpsrv.ListenAndServe()
}

func main() {
	// initialize variables
	var (
		serversInfo []interface{} = utils.LoadConfig("./config/servers.json")
		wg          sync.WaitGroup
	)

	// Initialize DDB tables
	actionspkg.InitializeTables(actionspkg.LoadDefaultConfig())

	// initialize servers in a go routine
	for _, srvInfo := range serversInfo {
		wg.Add(1)

		go func(srvInfo interface{}) {
			server := InitServer(srvInfo.(map[string]interface{}))
			StartServer(server)
			wg.Done()
		}(srvInfo)

	}

	// wait for all the servers to initialize before ending
	wg.Wait()
}
