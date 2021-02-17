package main

import (
	"log"

	"github.com/gorilla/websocket"
)

type responseRooms struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

func notifyPollPool(room string, newCount int) {
	response := responseRooms{room, newCount}
	for _, conn := range poll {
		if err := conn.WriteJSON(response); err != nil {
			log.Println(err)
		}
	}
}

func makeRoomList() []responseRooms {
	list := make([]responseRooms, len(rooms))
	next := 0
	for k, v := range rooms {
		list[next] = responseRooms{k, len(v)}
		next++
	}
	return list
}

func pollReader(conn *websocket.Conn) {
	// write all current rooms
	list := makeRoomList()
	if err := conn.WriteJSON(list); err != nil {
		log.Println(err)
	}
	// subscribe to the poll pool
	uuid, _ := makeUUID()
	poll[uuid] = conn
}
