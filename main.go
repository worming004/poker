package main

import (
	"flag"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type conf struct {
	password, port string
}

var passwordFlag *string = flag.String("password", "", "password to access")
var portFlag *string = flag.String("port", "8000", "port of host")

func main() {
	c := getConfig()
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

func getConfig() conf {
	flag.Parse()
	var password string
	passwordEnv := os.Getenv("PASSWORD")
	if len(passwordEnv) > 0 {
		password = passwordEnv
	} else {
		password = *passwordFlag
	}
	c := conf{
		password: password,
		port:     ":" + *portFlag,
	}
	return c
}
