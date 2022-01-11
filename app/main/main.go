package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
)

type Counter struct {
	Message string `json:"message"`
	Data    string `json:"data"`
}

var counterChannel chan uint64

type HTTPServer struct {
	http.Server
	shutdownReq chan bool
}

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
	server := &HTTPServer{
		Server: http.Server{
			Addr:    ":8000",
			Handler: router,
		},
		shutdownReq: make(chan bool),
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
			log.Printf("Listen and serve: %v", err)
		}
		done <- true
	}()

	//wait shutdown
	server.WaitShutdown()

	closeMain()

	<-done
	log.Printf("DONE!")
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

func (s *HTTPServer) WaitShutdown() {
	irqSig := make(chan os.Signal, 1)
	signal.Notify(irqSig, syscall.SIGINT, syscall.SIGTERM)

	//Wait interrupt or shutdown request through /shutdown
	select {
	case sig := <-irqSig:
		log.Printf("Shutdown request (signal: %v)", sig)
	case sig := <-s.shutdownReq:
		log.Printf("Shutdown request (/shutdown %v)", sig)
	}

	log.Printf("Stopping HTTP server ...")

	//Create shutdown context with 10 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//shutdown the server
	err := s.Shutdown(ctx)
	if err != nil {
		log.Printf("Shutdown request error: %v", err)
	}
}

func (s *HTTPServer) ShutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Shutdown server"))
	go func() {
		s.shutdownReq <- true
	}()
}
