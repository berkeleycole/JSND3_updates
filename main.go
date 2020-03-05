package main

import (
	"log"
	"net/http"
	"os"
	"time"
)

var dataFile string

func init() {
	dataFile = os.Getenv("DATA_FILE")
	if dataFile == "" {
		dataFile = "data.json"
	}
}

func main() {
	if err := generateTrackSegments(dataFile); err != nil {
		log.Fatal(err)
	}

	service, err := NewRaceService(DataFromJSONFile(dataFile))
	panicErr(err)

	server, err := NewServer(service)
	panicErr(err)

	srv := &http.Server{
		Handler:      server,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
