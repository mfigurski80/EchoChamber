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

var pool []*websocket.Conn = make([]*websocket.Conn, 0)

func reader(conn *websocket.Conn) {
	// subscribe to pool
	i := len(pool)
	pool = append(pool, conn)
	// write from websocket
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println(string(p))
		for next_i, c := range pool {
			if next_i == i {
				continue
			}
			go func(c *websocket.Conn) {
				if err := c.WriteMessage(messageType, p); err != nil {
					log.Println(err)
				}
			}(c)
		}
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
