package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var GIF = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

func main() {
	r := mux.NewRouter()
	// trackRepo := api.NewTrackPG("localhost:5432")
	// TODO
	// r.Host("www.example.com")

	// Add routes
	r.HandleFunc("/v1/event.gif", EventPixelHandler).Methods("GET")

	// Start server
	http.ListenAndServe(":8080", r)
}

func EventPixelHandler(w http.ResponseWriter, r *http.Request) {
	// Get pixel data from client
	v := r.URL.Query()
	rawEvent := v.Get("q")

	log.Println(rawEvent)

	// Write gif back to client
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Content-Type", "image/gif")
	w.Write(GIF)
}
