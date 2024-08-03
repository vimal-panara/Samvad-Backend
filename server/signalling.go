package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomId := AllRooms.CreateRoom()

	type resp struct {
		RoomId string `json:"room_id"`
	}

	fmt.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomId: roomId})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomId  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		msg := <-broadcast
		for _, client := range AllRooms.Map[msg.RoomId] {
			if client.Conn != msg.Client {
				err := client.Conn.WriteJSON(msg.Message)

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomId, ok := r.URL.Query()["roomId"]
	if !ok {
		log.Print("RoomId is missing in the parameter")
		return
	}
	log.Print(roomId)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomId[0], false, ws)

	go broadcaster()

	for {
		var msg broadcastMsg

		if err := ws.ReadJSON(&msg.Message); err != nil {
			log.Fatal("Read Error:", err)
		}

		msg.Client = ws
		msg.RoomId = roomId[0]

		log.Print(msg.Message)

		broadcast <- msg
	}
}

func DeleteRoomHandle(w http.ResponseWriter, r *http.Request) {
	roomId, ok := r.URL.Query()["roomId"]
	if !ok {
		log.Print("RoomId is missing in the parameter")
		return
	}
	log.Print(roomId)

	AllRooms.DeleteRoom(roomId[0])

	type resp struct {
		RoomId string `json:"room_id"`
	}

	fmt.Println(AllRooms.Map)

	json.NewEncoder(w).Encode(resp{RoomId: roomId[0]})
}
