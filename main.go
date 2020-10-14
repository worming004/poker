package main

import (
	"flag"
	"log"
	"time"
)

type conf struct {
	password, hostname string
}

var logState *bool = flag.Bool("logState", false, "log number of connection open")
var password *string = flag.String("password", "", "password to access")
var hostname *string = flag.String("hostname", "http://poker.craftlabit.be", "hostname of host")

func main() {
	flag.Parse()
	c := conf{
		password: *password,
		hostname: *hostname,
	}
	h := newHub(c)
	server := getApplicationServer(h, c)
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
