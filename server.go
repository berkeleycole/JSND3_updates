package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Server is responsible for exposing the services via HTTP
type Server struct {
	router  http.Handler
	service *RaceService
}

// NewServer returns a new http server
func NewServer(service *RaceService) (*Server, error) {
	router := makeRouter(service)

	server := Server{
		router:  router,
		service: service,
	}

	return &server, nil
}

// ServerOpt defines an option that can be applied to a server
// to help configure it.

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func makeRouter(service *RaceService) http.Handler {
	router := mux.NewRouter()
	r := router.PathPrefix("/api/").Subrouter()

	// cars
	r.HandleFunc("/cars", listCars(service)).Methods("GET")

	// tracks
	r.HandleFunc("/tracks", listTracks(service)).Methods("GET")

	// races
	r.HandleFunc("/races", unimplemented).Methods("GET")
	r.HandleFunc("/races", unimplemented).Methods("POST")
	r.HandleFunc("/races/{raceID}", unimplemented).Methods("GET")

	r.HandleFunc("/", notFound)
	router.HandleFunc("/", notFound)

	return router
}

func listCars(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(service.Cars)
		panicErr(err)
	})
}

func listTracks(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(service.Tracks)
		panicErr(err)
	})
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
}

func unimplemented(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(501)
}

func panicErr(err error) {
	if err != nil {
		panic(err)
	}
}
