package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// holds connection
var clients = make(map[*websocket.Conn]bool)

// allows us to broad message on ws to all connected clients
var broadcast = make(chan Job)

// Configure the upgrader
var upgrader = websocket.Upgrader{}

// this handles websocket connection
func handleConnections(w http.ResponseWriter, r *http.Request) {
	log.Printf("We got a new websocket connection from %s\n", r.Header.Get("Origin"))
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	// Register our new client
	clients[conn] = true
	for {
		msg := "{}"
		// Read in a new message as JSON and map it to a Message object
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, conn)
			break
		}
		// Send the newly received message to the broadcast channel
		for {
			msg := "{}"
			// Read in a new message as JSON and map it to a Message object
			err := conn.ReadJSON(&msg)
			if err != nil {
				log.Printf("error: %v", err)
				// let's remove connection from map
				delete(clients, conn)
				break
			}
		}
	}
	// Make sure we close the connection when the function returns
	defer conn.Close()
}

func StartServer() {
	//var port = flag.Int("port", 8080, "Specify port for HTTP Interface, default : 8080")
	flag.Parse()
	port := 8080
	routes := mux.NewRouter()
	routes.PathPrefix("/").Handler(http.FileServer(http.Dir("./public/")))

	// send message from broadcast channel to all connected clients
	go handleMessages()
	// setup ws route
	http.HandleFunc("/ws", handleConnections)
	// setup our REST endpoints
	http.HandleFunc("/monitor/job/add", handleAddJob)
	http.HandleFunc("/monitor/job/remove", handleDeleteJob)
	http.Handle("/", routes)

	log.Printf("Started HTTP Interface on port %d\n", port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func handleAddJob(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var reqJob Job

		err := json.Unmarshal(body, &reqJob)
		if err != nil {
			http.Error(w, "Bad Request: Failed to parse JSON", http.StatusBadRequest)
		}
		// Create a new record.
		if _, ok := monitoringJobs[reqJob.URL]; !ok {
			go Monitor(reqJob)
		}
		monitoringJobs[reqJob.URL] = reqJob
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

	} else {
		http.Error(w, "Bad Request: Method not allowed", http.StatusBadRequest)
	}
}

func handleDeleteJob(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		body, _ := ioutil.ReadAll(r.Body)
		var reqJob Job

		err := json.Unmarshal(body, &reqJob)
		if err != nil {
			http.Error(w, "Bad Request: Failed to parse JSON", http.StatusBadRequest)
		}
		// delete record.
		delete(monitoringJobs, reqJob.URL)
	} else {
		http.Error(w, "Bad Request:  Method not allowed", http.StatusBadRequest)
	}
}
