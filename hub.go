package main

import (
	"context"
	"encoding/json"
	"errors"
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

func (pss PlayerSelectionState) Validate() (bool, error) {
	if len(pss.Player) > 35 {
		return false, errors.New("player name too long")
	}
	return true, nil
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
	// happens when user quit before connecting to a room
	if room == nil {
		return nil
	}
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

func (h *hub) handleCommand(ctx context.Context, cmd ExternalCommand, currentClient *virtualClient) {
	realCommand := commandFactory(ctx, cmd, currentClient)
	getLogger(ctx).Infof("command received %T: %+v", realCommand, realCommand)

	if cmd, ok := realCommand.(addPlayerCommand); ok {
		currentClient.roomID = cmd.roomID
		currentClient.playerID = cmd.playerID
	}
	room, ok := h.room[cmd.RoomID]
	if !ok {
		room = newRoom(cmd.RoomID)
		h.room[cmd.RoomID] = room
		time.AfterFunc(4*time.Hour, func() { h.deleteRoom(ctx, cmd.RoomID) })
	}
	room.Connections[currentClient] = currentClient.Conn
	err := realCommand.Do(ctx, room)
	if err != nil {
		return
	}
	room.broadcastCurrentState()
}

func (h *hub) handleSocket(w http.ResponseWriter, r *http.Request) {
	password := r.URL.Query()["password"][0]
	ctx := r.Context()
	log := getLogger(ctx)
	log.Println("Attempt to connect")
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
		var cmd ExternalCommand
		_, message, err := c.ReadMessage()
		ctx = pushNewSendID(ctx)
		log = getLogger(ctx)
		if err != nil {
			h.disconnect(ctx, client)
			log.Printf("close ws because of %s\n", err)
			break
		}
		err = json.Unmarshal(message, &cmd)
		if err != nil {
			panic(err)
		}
		h.handleCommand(ctx, cmd, client)
	}
}

func (h *hub) disconnect(ctx context.Context, c *virtualClient) error {
	delete(h.clients, c)
	room, ok := h.room[c.roomID]
	if ok {
		delete(room.Players, c.playerID)
		delete(room.Connections, c)
		if len(room.Connections) == 0 {
			h.deleteRoom(ctx, c.roomID)
		}
	}
	room.broadcastCurrentState()
	return nil
}

func (h *hub) deleteRoom(ctx context.Context, roomID int) {
	if p, ok := h.room[roomID]; ok {
		getLogger(ctx).Printf("Deleting roomID %d\n", roomID)
		for c := range p.Connections {
			c.Close()
		}
		delete(h.room, roomID)
	} else {
		getLogger(ctx).Printf("Cannot delete roomID %d because it's not found\n", roomID)
	}
}
