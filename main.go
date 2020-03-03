package main

import (
	"net/http"
	"os"
	"time"
)

var dataFile string

func init() {
	os.Getenv("DATA_FILE")
}

func main() {
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
