package main

import (
	"net/http"
)

func getApplicationServer(h hub) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleSocket)
	return &http.Server{
		Addr:    "localhost:6000",
		Handler: mux,
	}
}
