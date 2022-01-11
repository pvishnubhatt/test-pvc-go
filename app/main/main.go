package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Counter struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

var counterChannel chan uint64

func initMain() {
	numChannels := 16
	counterChannel = make(chan uint64, numChannels)
}

func main() {
	initMain()
	router := mux.NewRouter()
	router.HandleFunc("/", handleMain)
	router.HandleFunc("/counter", handleMain)
	router.HandleFunc("/counter/get", handleMainGet)

	log.Println("Main Server is running!")
	fmt.Println(http.ListenAndServe(":8000", router))
}

func handleMain(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleMain")
	response := map[string]string{
		"message": "Welcome to test-pvc - Main",
	}
	json.NewEncoder(rw).Encode(response)
}

func handleMainGet(rw http.ResponseWriter, r *http.Request) {
	log.Println("main.handleMainGet")
	go getCounterFromChannel(counterChannel)
	counter := <-counterChannel
	response := map[string]string{
		"message": "Welcome to test-pvc - Get",
		"counter": strconv.FormatUint(counter, 10),
	}
	json.NewEncoder(rw).Encode(response)
}

func getCounterFromChannel(retChannel chan uint64) {
	log.Println("main.getCounterFromChannel ", cap(retChannel), len(retChannel))
	retChannel <- getCounter()
}

func getCounter() uint64 {
	log.Println("main.getCounter")
	resp, err := http.Get("http://get-go-service:8000/counter/get/get")
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var counter Counter
	json.Unmarshal(bodyBytes, &counter)
	ctr, err := strconv.ParseUint(counter.Data, 10, 64)
	if err != nil {
		log.Fatalln(err)
	}
	return ctr
}
