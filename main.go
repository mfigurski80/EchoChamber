package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var pool chan []byte = make(chan []byte)

func reader(conn *websocket.Conn) {
	// go read from channel
	go func() {
		for v := range pool {
			conn.WriteMessage(1, v)
		}
	}()
	// write from websocket
	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		pool <- p

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected...")
	reader(ws)
}

func main() {
	http.HandleFunc("/ws", wsEndpoint)
	log.Println("Serving on http://0.0.0.0:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
