package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/lib/pq"

	track "github.com/diericx/tracker/backend/pkg"
	"github.com/gorilla/mux"
)

var gif = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

func main() {
	trackPostgresRepo, err := track.NewPostgresRepo("localhost", 5432, "postgres", "postgres", "mattribution")
	if err != nil {
		panic(err)
	}
	defer trackPostgresRepo.GetDB().Close()
	trackService := track.NewService(trackPostgresRepo)

	r := mux.NewRouter()
	// trackRepo := api.NewTrackPG("localhost:5432")
	// TODO
	// r.Host("www.example.com")

	// Add routes
	r.HandleFunc("/v1/pixel/track", func(w http.ResponseWriter, r *http.Request) {
		// Get pixel data from client
		v := r.URL.Query()
		rawEvent := v.Get("data")
		data, err := base64.StdEncoding.DecodeString(rawEvent)
		if err != nil {
			panic(err)
			return
		}

		log.Println(data)

		track := track.Track{}
		if err := json.Unmarshal(data, &track); err != nil {
			panic(err)
		}

		log.Println(track)

		id, err := trackService.New(track)
		if err != nil {
			panic(err)
		}

		log.Println(id)

		// Write gif back to client
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gif)
	}).Methods("GET")

	// Start server
	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}
