package main

import (
	"encoding/json"
	"log"
	"net/http"
	"test-pvc-go/test-pvc-go/app/common"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	server := &common.HTTPServer{
		Server: http.Server{
			Addr:    ":8000",
			Handler: router,
		},
		ShutdownReq: make(chan bool),
	}
	router.HandleFunc("/counter/get/get", handleGet)

	log.Println("Get Server is running!")
	done := make(chan bool)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			log.Printf("Main Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	<-done
	log.Printf("Get DONE!")
}

func handleGet(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleGet")
	response := map[string]string{
		"message": "Welcome to test-pvc - Get",
		"data":    "29",
	}
	json.NewEncoder(rw).Encode(response)
}
