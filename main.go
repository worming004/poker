package main

import (
	"flag"
	"time"

	"github.com/sirupsen/logrus"
)

type conf struct {
	password, port string
}

var password *string = flag.String("password", "", "password to access")
var port *string = flag.String("port", "8000", "port of host")

func main() {
	flag.Parse()
	c := conf{
		password: *password,
		port:     ":" + *port,
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
