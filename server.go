package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func getApplicationServer(h *hub, c conf) *http.Server {
	mux := http.NewServeMux()
	indexHandler := getIndexHandler(Model{Hostname: c.hostname})
	mux.Handle("/static/", getStaticHandler("/static/"))
	mux.HandleFunc("/connect", h.handleSocket)
	mux.HandleFunc("/cards", getCardHandler(h))
	mux.HandleFunc("/newid", getNewIDHandler())
	mux.HandleFunc("/", indexHandler)
	return &http.Server{
		Addr:    "localhost:6000",
		Handler: mux,
	}
}

func getStaticHandler(prefix string) http.Handler {
	fs := http.FileServer(http.Dir("./static"))
	if prefix != "" {
		return fs
	}
	return http.StripPrefix(prefix, fs)
}

type Model struct {
	Hostname string
}

func getIndexHandler(m Model) http.HandlerFunc {
	tmpl := template.Must(template.ParseGlob("./html/index.gohtml"))
	return func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, m)
		if err != nil {
			log.Panic(err)
		}
	}
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
