package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"
)

type conf struct {
	password, hostname, port string
}

var password *string = flag.String("password", "", "password to access")
var hostname *string = flag.String("hostname", "http://poker.craftlabit.be", "hostname of host")
var port *string = flag.String("port", "6000", "port of host")

func main() {
	flag.Parse()
	c := conf{
		password: *password,
		hostname: *hostname,
		port:     *port,
	}
	h := newHub(c)
	server := getApplicationServer(h, c)

	go func() {
		ticker := time.NewTicker(12 * time.Hour)
		logger := logrus.New()
		for {
			<-ticker.C
			logNumberOfCurrentSession(logger, h)
		}
	}()

	logrus.Info("Starting server")
	logrus.Panic(server.ListenAndServe())
}
