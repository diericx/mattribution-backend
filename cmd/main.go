package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net"
	"net/http"

	_ "github.com/lib/pq"

	track "github.com/diericx/tracker/backend/pkg"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var gif = []byte{
	71, 73, 70, 56, 57, 97, 1, 0, 1, 0, 128, 0, 0, 0, 0, 0,
	255, 255, 255, 33, 249, 4, 1, 0, 0, 0, 0, 44, 0, 0, 0, 0,
	1, 0, 1, 0, 0, 2, 1, 68, 0, 59,
}

const (
	mockOwnerID = 1
)

type tracksPayload struct {
	Tracks []track.Track `json:"tracks"`
}
type trackDailyCountsPayload struct {
	DailyCounts []track.DailyCount `json:"dailyCounts"`
}

func main() {
	trackPostgresRepo, err := track.NewPostgresRepo("localhost", 5432, "postgres", "postgres", "mattribution")
	if err != nil {
		panic(err)
	}
	defer trackPostgresRepo.GetDB().Close()
	trackService := track.NewService(trackPostgresRepo)

	r := mux.NewRouter()
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
		}

		track := track.Track{}
		if err := json.Unmarshal(data, &track); err != nil {
			panic(err)
		}
		ip, _, _ := net.SplitHostPort(r.RemoteAddr)
		track.IP = ip

		_, err = trackService.New(mockOwnerID, track)
		if err != nil {
			panic(err)
		}

		// Write gif back to client
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Content-Type", "image/gif")
		w.Write(gif)
	}).Methods("GET")

	r.HandleFunc("/v1/pixel/tracks", func(w http.ResponseWriter, r *http.Request) {
		// Get pixel data from client
		// v := r.URL.Query()

		tracks, err := trackService.GetAll(mockOwnerID)
		if err != nil {
			panic(err)
		}

		p := tracksPayload{tracks}

		// Write gif back to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}).Methods("GET")

	r.HandleFunc("/v1/pixel/tracks/count", func(w http.ResponseWriter, r *http.Request) {
		// Get pixel data from client
		// v := r.URL.Query()

		trackCounts, err := trackService.GetDailyCounts(mockOwnerID)
		if err != nil {
			panic(err)
		}

		p := trackDailyCountsPayload{trackCounts}

		// Write gif back to client
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(p)
	}).Methods("GET")

	// Start server
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(r)))

}
