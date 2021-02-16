package main

import (
	"log"

	"github.com/gorilla/websocket"
)

func roomReader(rm string, conn *websocket.Conn) {
	// subscribe to pool
	uuid, _ := makeUUID()
	rooms[rm][uuid] = conn
	// notify poll pool
	notifyPollPool(rm, len(rooms[rm]))
	// write to pool from websocket
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break // connection closed
		}
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
