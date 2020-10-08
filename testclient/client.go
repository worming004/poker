package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:6000", "http service address")

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	// defer c.Close()

	message := "{\"player\":\"worming\"}"

	c.WriteMessage(websocket.TextMessage, []byte(message))
	fmt.Println("closing")
	time.Sleep(4 * time.Second)

	err = c.Close()
	if err != nil {
		panic(err)
	}
	fmt.Println("closed")
	time.Sleep(4 * time.Second)

}
