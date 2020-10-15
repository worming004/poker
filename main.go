package main

import (
	"flag"

	"github.com/sirupsen/logrus"
)

type conf struct {
	password, hostname string
}

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

	logrus.Info("Starting server")
	logrus.Panic(server.ListenAndServe())
}
