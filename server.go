package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

// Server is responsible for exposing the services via HTTP
type Server struct {
	router http.Handler
}

// NewServer returns a new http server
func NewServer() (*Server, error) {
	router := makeRouter()
	server := Server{
		router: router,
	}

	return &server, nil
}

// ServerOpt defines an option that can be applied to a server
// to help configure it.

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func makeRouter() http.Handler {
	router := mux.NewRouter()
	r := router.PathPrefix("/api/").Subrouter()

	// cars
	r.HandleFunc("/cars", listCars).Methods("GET")

	// tracks
	r.HandleFunc("/tracks", listTracks).Methods("GET")

	// races
	r.HandleFunc("/races", unimplemented).Methods("GET")
	r.HandleFunc("/races", unimplemented).Methods("POST")
	r.HandleFunc("/races/{raceID}", unimplemented).Methods("GET")

	r.HandleFunc("/", notFound)
	router.HandleFunc("/", notFound)

	return router
}

func listCars(w http.ResponseWriter, r *http.Request) {
	cars := readFile().Cars
	err := json.NewEncoder(w).Encode(cars)
	panicErr(err)
}

func listTracks(w http.ResponseWriter, r *http.Request) {
	tracks := readFile().Tracks
	err := json.NewEncoder(w).Encode(tracks)
	panicErr(err)
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func unimplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
}

type dataJSON struct {
	Cars   []*Car   `json:"cars"`
	Tracks []*Track `json:"tracks"`
}

func readFile() dataJSON {
	f, err := os.Open("./data.json")
	panicErr(err)

	var data dataJSON

	err = json.NewDecoder(f).Decode(&data)
	panicErr(err)

	return data
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
