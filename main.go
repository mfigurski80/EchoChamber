package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Pool []*websocket.Conn

var rooms map[string]Pool = make(map[string]Pool, 0)

func reader(rm string, conn *websocket.Conn) {
	// subscribe to pool
	i := len(rooms[rm])
	rooms[rm] = append(rooms[rm], conn)
	// write from websocket
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break // this where checking close?
		}
		log.Println(string(p))
		for next_i, c := range rooms[rm] {
			if next_i == i {
				continue
			}
			if err := c.WriteMessage(messageType, p); err != nil {
				log.Println(err)
			}
		}
	}
	// remove from pool
	rooms[rm] = append(rooms[rm][:i], rooms[rm][i+1:]...)
	if len(rooms[rm]) == 0 {
		delete(rooms, rm)
	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	var room string
	queryRoom, ok := r.URL.Query()["room"]
	if !ok || len(queryRoom) == 0 {
		room = ""
	}
	room = queryRoom[0]

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
	http.HandleFunc("/ws", wsEndpoint)
	log.Println("Serving on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
