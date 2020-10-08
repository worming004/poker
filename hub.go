package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type card struct{}

type player struct {
	name string
}

type state string

type playerSelectionState struct {
	card
	player
	state
}

type party struct {
	players []playerSelectionState
}

type hub struct {
	upgrader websocket.Upgrader
	party
}

func newHub() *hub {
	return new(hub)
}

func (h *hub) handleCommand(cmd command) {
	realCommand := commandFactory(cmd)
	realCommand.Do(h)
}

func (h *hub) handleSocket(w http.ResponseWriter, r *http.Request) {
	c, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(err)
	}
	defer c.Close()
	for {
		var cmd command
		err := c.ReadJSON(&cmd)
		if err != nil {
			log.Panic(err)
		}
		h.handleCommand(cmd)
	}
}
