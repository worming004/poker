package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

func getApplicationServer(h *hub) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", serveStatic)
	mux.HandleFunc("/static", serveStatic)
	mux.HandleFunc("/connect", h.handleSocket)
	mux.HandleFunc("/cards", getCardHandler(h))
	mux.HandleFunc("/newid", getNewIDHandler())
	return &http.Server{
		Addr:    "localhost:6000",
		Handler: mux,
	}
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	fs := http.FileServer(http.Dir("./static"))
	log.Println("serving static file")
	fs.ServeHTTP(w, r)
}

func getCardHandler(h *hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID, err := strconv.Atoi(r.URL.Query()["roomid"][0])
		if err != nil {
			panic(err)
		}
		msg := struct {
			Cards []Card
		}{
			Cards: h.room[roomID].SetOfCards,
		}
		b, err := json.Marshal(msg)
		if err != nil {
			panic(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(b)
	}
}

type getNewID struct {
	ID int `json:"id"`
}

func getNewIDHandler() http.HandlerFunc {
	current := 0
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("serving id")
		current++
		res := getNewID{current}
		b, err := json.Marshal(res)
		if err != nil {
			log.Printf("err while parsing res %s\n", err)
		}
		w.Write(b)
	}
}
