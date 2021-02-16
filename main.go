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

type Room map[string]*websocket.Conn

var rooms map[string]Room = make(map[string]Room, 0)

func reader(rm string, conn *websocket.Conn) {
	// subscribe to pool
	uuid, _ := makeUUID()
	rooms[rm][uuid] = conn
	// write to pool from websocket
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break // connection closed
		}
		log.Println(string(p))
		for k, c := range rooms[rm] {
			if k == uuid {
				continue
			}
			if err := c.WriteMessage(messageType, p); err != nil {
				log.Println(err)
			}
		}
	}
	// remove from pool
	delete(rooms[rm], uuid)
	if len(rooms[rm]) == 0 {
		delete(rooms, rm)
	}
}

type ResponseRooms struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func pollEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	list := make([]ResponseRooms, len(rooms))
	next := 0
	for k, v := range rooms {
		list[next] = ResponseRooms{k, len(v)}
		next++
	}
	json.NewEncoder(w).Encode(list)
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

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

	log.Println("client connected")
	reader(room, ws)
}

func mainEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello! This is a websocket echo server for Rohan's project"))
}

func main() {
	dir, _ := os.Getwd()
	log.Println(dir)
	http.HandleFunc("/", mainEndpoint)
	http.HandleFunc("/poll", pollEndpoint)
	http.HandleFunc("/ws", wsEndpoint)
	log.Println("Serving on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
