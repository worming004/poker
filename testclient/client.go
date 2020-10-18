package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

type Command struct {
	Player  string
	Action  string
	Payload map[string]string
}

var addr = flag.String("addr", "localhost:8080", "http service address")
var playerName = flag.String("pn", "worming", "player name")

func main() {
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/connect"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		_, msg, err := c.ReadMessage()
		if err != nil {
			fmt.Printf("error while reading. err : %v\n", err)
		}
		fmt.Printf("message received: %s\n", string(msg))
	}()

	cmd := Command{
		Player: *playerName,
		// Action: "actionAddPlayer",
		Action: "actionRoomShowCard",
		// Action:  "actionRoomResetCard",
		Payload: map[string]string{},
	}

	message, err := json.Marshal(cmd)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sending %v\n", message)

	c.WriteMessage(websocket.TextMessage, []byte(message))

	time.Sleep(2 * time.Second)

	err = c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		panic(err)
	}
}
