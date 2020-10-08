package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type PlayerState string

const (
	playerSelectedCard   = "playerSelected"
	playerNoSelectedCard = "playerNoSelectedCard"
)

type PlayerSelectionState struct {
	Player string
	PlayerState
	Card string
}

func newPlayerSelectionState(playerName string) *PlayerSelectionState {
	return &PlayerSelectionState{Card: "", Player: playerName, PlayerState: playerNoSelectedCard}
}

type RoomState string

const (
	roomSelecting = "roomSelecting"
	roomShow      = "roomShow"
)

type room struct {
	RoomID      int
	Players     map[int]*PlayerSelectionState
	Connections map[*virtualClient]*websocket.Conn
	RoomState   RoomState
	SetOfCards  Cards
}

type roomResponse struct {
	Players   map[int]*PlayerSelectionState
	RoomState RoomState
}

func newRoom(roomID int) *room {
	return &room{
		RoomID:      roomID,
		RoomState:   roomSelecting,
		SetOfCards:  scrumCards,
		Connections: make(map[*virtualClient]*websocket.Conn),
		Players:     make(map[int]*PlayerSelectionState),
	}
}

type hub struct {
	clients  map[*virtualClient]bool
	room     map[int]*room
	password string
}

type ServerAction struct {
	ActionType string      `json:"actionType"`
	Payload    interface{} `json:"payload"`
}

func (room *room) broadcastCurrentState() error {
	for _, client := range room.Connections {
		client.WriteJSON(ServerAction{
			ActionType: "refreshState",
			Payload: roomResponse{
				Players:   room.Players,
				RoomState: room.RoomState,
			},
		})
	}
	return nil
}

func newHub(c conf) *hub {
	h := new(hub)
	h.clients = make(map[*virtualClient]bool)
	h.room = make(map[int]*room)
	h.password = c.password
	return h
}

func (h *hub) handleCommand(cmd Command, currentClient *virtualClient) {
	realCommand := commandFactory(cmd, currentClient)
	if cmd, ok := realCommand.(addPlayerCommand); ok {
		currentClient.roomID = cmd.roomID
		currentClient.playerID = cmd.playerID
	}
	room, ok := h.room[cmd.RoomID]
	if !ok {
		room = newRoom(cmd.RoomID)
		h.room[cmd.RoomID] = room
		time.AfterFunc(4*time.Hour, func() { h.deleteRoom(cmd.RoomID) })
	}
	room.Connections[currentClient] = currentClient.Conn
	realCommand.Do(room)
	room.broadcastCurrentState()
}

func (h *hub) handleSocket(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query()["password"][0]
	if password != h.password {
		log.Println("socket invalid due to wrong password")
		return
	}
	playerid, err := strconv.Atoi(r.URL.Query()["playerid"][0])
	if err != nil {
		log.Printf("not able to parse %s as int\n", r.URL.Query()["playerid"][0])
		return
	}
	if playerid == 0 {
		log.Printf("invalid playerID %d\n", playerid)
		return
	}
	log.Println("opening socket")
	conn := websocket.Upgrader{}
	c, err := conn.Upgrade(w, r, nil)
	if err != nil {
		log.Panic(err)
	}
	client := &virtualClient{h, c, 0, 0}
	h.clients[client] = true
	defer c.Close()
	for {
		var cmd Command
		_, message, err := c.ReadMessage()
		if err != nil {
			h.disconnect(client)
			log.Printf("close ws because of %s\n", err)
			break
		}
		log.Printf("received : %s", string(message))
		err = json.Unmarshal(message, &cmd)
		if err != nil {
			panic(err)
		}
		h.handleCommand(cmd, client)
	}
}

func (h *hub) disconnect(c *virtualClient) error {
	delete(h.clients, c)
	room := h.room[c.roomID]
	delete(room.Players, c.playerID)
	delete(room.Connections, c)
	room.broadcastCurrentState()
	return nil
}

func (h *hub) deleteRoom(roomID int) {
	p := h.room[roomID]
	for c := range p.Connections {
		c.Close()
	}
	delete(h.room, roomID)
}

func (h *hub) logState() {
	log.Printf("number of open connection : %d", len(h.clients))
}
