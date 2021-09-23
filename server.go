package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"

	"github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var defaultLogger *logrus.Logger

func getApplicationServer(h *hub, c conf) *http.Server {
	mux := mux.NewRouter()
	indexHandler := getIndexHandler(Model{Hostname: c.hostname})
	mux.HandleFunc("/", indexHandler).Methods("GET")
	// Use extra {route} in order to allow file discovery. https://stackoverflow.com/questions/21234639/golang-gorilla-mux-with-http-fileserver-returning-404
	mux.Handle("/static/{route}", getStaticHandler("/static/")).Methods("GET")
	mux.HandleFunc("/connect", h.handleSocket).Methods("GET")
	mux.HandleFunc("/cards", getCardHandler(h)).Methods("GET")
	mux.HandleFunc("/newid", getNewIDHandler()).Methods("GET")
	mux.Use(loggerMiddleware, allowCors)
	return &http.Server{
		Addr:    "localhost:" + c.port,
		Handler: mux,
	}
}

func allowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func getStaticHandler(prefix string) http.Handler {
	fs := http.FileServer(http.Dir("./static"))
	if prefix == "" {
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
			logrus.Panic(err)
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
		logrus.Println("serving id")
		current++
		if current == 4294967294 {
			current = 0
		}
		res := getNewID{current}
		b, err := json.Marshal(res)
		if err != nil {
			logrus.Printf("err while parsing res %s\n", err)
		}
		w.Write(b)
	}
}
