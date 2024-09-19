package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/rs/cors"
	"github.com/vp-0312/Samvad-Backend/server"
)

func main() {
	fmt.Println("Welcome to the Samvad!")

	server.AllRooms.Init()

	mux := http.NewServeMux()

	// Set up CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace with your Remix frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},                                                                            // Adjust allowed methods
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	c.Handler(mux)

	mux.HandleFunc("/create", server.CreateRoomRequestHandler)
	mux.HandleFunc("/join", server.JoinRoomRequestHandler)
	mux.HandleFunc("/delete", server.DeleteRoomHandle)

	log.Print("starting server on port 8000")

	if err := http.ListenAndServe(":8000", nil); err != nil {
		log.Fatal("can not start server:", err)
	}

}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("setting cors")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		if r.Method == http.MethodOptions {
			// Allow preflight requests and return OK without processing further
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
