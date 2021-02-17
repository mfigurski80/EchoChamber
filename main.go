package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type room map[string]*websocket.Conn

var rooms = make(map[string]room, 0)

var poll = make(room, 0)

func pollEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(makeRoomList())
}

func wsPollEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	pollReader(ws)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// get room url query
	var room string
	queryRoom, ok := r.URL.Query()["room"]
	if !ok || len(queryRoom) == 0 {
		room = ""
	} else {
		room = queryRoom[0]
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Printf("client connected to room '%s'\n", room)
	roomReader(room, ws)
}

func mainEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello! This is a websocket echo server for Rohan's project"))
}

func main() {
	dir, _ := os.Getwd()
	log.Println(dir)
	http.HandleFunc("/", mainEndpoint)
	http.HandleFunc("/poll", pollEndpoint)
	http.HandleFunc("/poll/ws", wsPollEndpoint)
	http.HandleFunc("/ws", wsEndpoint)
	log.Println("Serving on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
