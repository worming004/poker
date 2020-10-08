package main

import "github.com/gorilla/websocket"

// virtualClient is a hub client hosted on Server. One by each ws connection
type virtualClient struct {
	*hub
	*websocket.Conn
	roomID   int
	playerID int
}
