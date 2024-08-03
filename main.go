package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/vp-0312/Samvad-Backend/server"
)

func main() {
	fmt.Println("Welcome to the Samvad!")

	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/delete", server.DeleteRoomHandle)

	log.Print("starting server on port 8000")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("can not start server:", err)
	}

}
