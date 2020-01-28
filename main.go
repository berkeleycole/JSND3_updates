package main

import (
	"net/http"
	"time"
)

func main() {
	server, err := NewServer()
	panicErr(err)

	srv := &http.Server{
		Handler:      server,
		Addr:         "0.0.0.0:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}
