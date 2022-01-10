package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/counter/get/get", handleGet)

	log.Println("Get Server is running!")
	fmt.Println(http.ListenAndServe(":8000", router))
}

func handleGet(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleGet")
	response := map[string]string{
		"message": "Welcome to test-pvc - Get",
		"data":    "29",
	}
	json.NewEncoder(rw).Encode(response)
}
