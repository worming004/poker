package main

import (
	"flag"
	"log"
	"time"
)

type conf struct {
	password string
}

var logState *bool = flag.Bool("logState", false, "log number of connection open")
var password *string = flag.String("password", "", "log number of connection open")

func main() {
	flag.Parse()
	c := conf{
		password: *password,
	}
	h := newHub(c)
	server := getApplicationServer(h)
	log.Println("starting server")

	if *logState {
		go func() {
			ticker := time.NewTicker(2 * time.Second)
			for {
				<-ticker.C
				h.logState()
			}
		}()
	}

	log.Panic(server.ListenAndServe())
}
