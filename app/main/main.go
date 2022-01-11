package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"test-pvc-go/test-pvc-go/app/common"

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

func closeMain() {
	close(counterChannel)
}

func main() {
	initMain()

	router := mux.NewRouter()
	server := &common.HTTPServer{
		Server: http.Server{
			Addr:    ":8000",
			Handler: router,
		},
		ShutdownReq: make(chan bool),
	}

	router.HandleFunc("/", handleMain)
	router.HandleFunc("/counter", handleMain)
	router.HandleFunc("/counter/get", handleMain)
	router.HandleFunc("/shutdown", server.ShutdownHandler)
	log.Println("Main Server is running!")
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

	closeMain()

	<-done
	log.Printf("Main DONE!")
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
	retChannel <- getCounter()
}

func getCounter() uint64 {
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
