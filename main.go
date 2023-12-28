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

func StartServer(chatServer *models.Server, wg *sync.WaitGroup) {
	defer wg.Done()

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
	var (
		wg sync.WaitGroup
		srv1 *models.Server = &models.Server{}
		srv2 *models.Server = &models.Server{}
	)
		
	// setup registry
	utils.InitRegistry()
	
	// setup servers and their respective message queues
	srv1.SetupServer("localhost:8080")
	wg.Add(1)
	go StartServer(srv1, &wg)

	srv2.SetupServer("localhost:8081")
	wg.Add(1)
	go StartServer(srv2, &wg)

	// wait for all the servers to start 
	wg.Wait()
}