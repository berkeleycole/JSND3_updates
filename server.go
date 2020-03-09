package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
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
	// if r.Method == http.MethodOptions {
	w.Header().Set("Access-Control-Allow-Origin", "http://lvh.me:3000")
	// w.Header().Set("Content-Type", "application/json")
	// 	w.Header().Add("Access-Control-Allow-Headers", "*")
	// 	w.WriteHeader(200)
	// 	return
	// }

	s.router.ServeHTTP(w, r)

	// fmt.Printf("+%v", r)
}

func makeRouter(service *RaceService) http.Handler {
	router := mux.NewRouter()
	r := router.PathPrefix("/api/").Subrouter()

	// cars
	r.HandleFunc("/cars", listCars(service)).Methods("GET")

	// tracks
	r.HandleFunc("/tracks", listTracks(service)).Methods("GET")

	// races
	r.HandleFunc("/races", listRaces(service)).Methods("GET")
	r.HandleFunc("/races", createRace(service)).Methods("POST")
	r.HandleFunc("/races/{raceID}", getRace(service)).Methods("GET")
	r.HandleFunc("/races/{raceID}/accelerate", accelerate(service)).Methods("POST")

	r.HandleFunc("/", notFound)
	router.HandleFunc("/", notFound)

	headersOk := handlers.AllowedHeaders([]string{"Origin", "Accept", "X-Requested-With", "Content-Type", "Access-Control-Allow-Origin"})
	// originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	return handlers.CORS(headersOk, methodsOk)(router)
	// return router
}

func accelerate(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		raceID, err := strconv.ParseUint(params["raceID"], 10, 64)
		if err != nil {
			panicErr(err)
		}

		err = service.Accelerate(uint(raceID))
		if err != nil {
			panicErr(err)
			return
		}
	})
}

func getRace(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		raceID, err := strconv.ParseUint(params["raceID"], 10, 64)
		if err != nil {
			panicErr(err)
		}

		race, err := service.GetRace(uint(raceID))
		if err != nil {
			panicErr(err)
			return
		}

		err = json.NewEncoder(w).Encode(race.Results)
		if err != nil {
			panicErr(err)
			return
		}
	})
}

type createRaceParams struct {
	PlayerID int `json:"player_id"`
	TrackID  int `json:"tack_id"`
}

func createRace(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var params createRaceParams

		if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
			panicErr(err)
			return
		}

		race, err := service.CreateRace(params.PlayerID, params.TrackID)
		if err != nil {
			panicErr(err)
			return
		}

		err = json.NewEncoder(w).Encode(race)
		panicErr(err)

		// go func() {
		// 	time.Sleep(3)
		// 	race.Start()
		// }()
	})
}

func thing() {
}

func listRaces(service *RaceService) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := json.NewEncoder(w).Encode(service.Races)
		panicErr(err)
	})
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
		log.Printf("got an error: %v", err)
		panic(err)
	}
}
