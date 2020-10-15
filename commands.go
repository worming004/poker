package main

import (
	"context"
	"fmt"
)

type action string

const (
	actionAddPlayer        = "actionAddPlayer"
	actionPlayerSelectCard = "actionPlayerSelectCard"
	actionRoomShowCard     = "actionRoomShowCard"
	actionRoomResetCard    = "actionRoomResetCard"
)

type invalidCommand struct{}

func (invcmd invalidCommand) Error() string {
	return "Invalid commande"
}

type playerNotFoundError struct {
	playerID int
}

func (p playerNotFoundError) Error() string {
	return fmt.Sprintf("player %d not found", p.playerID)
}

type Command struct {
	PlayerID int
	RoomID   int
	Action   action
	Payload  map[string]string
}

type funcCommand interface {
	Do(context.Context, *room) error
}

type addPlayerCommand struct {
	playerName string
	playerID   int
	roomID     int
	client     *virtualClient
}

func (cmd addPlayerCommand) Do(ctx context.Context, p *room) error {
	p.Players[cmd.playerID] = newPlayerSelectionState(cmd.playerName)
	if cmd.client != nil && cmd.client.Conn != nil {
		cmd.client.Conn.WriteJSON(ServerAction{
			ActionType: "refreshCards",
			Payload: struct {
				Cards []Card `json:"cards"`
			}{
				Cards: p.SetOfCards,
			},
		})
	}
	return nil
}

type playerSelectCardCommand struct {
	playerID int
	card     Card
}

func (cmd playerSelectCardCommand) Do(ctx context.Context, p *room) error {
	if !p.SetOfCards.contains(cmd.card) {
		return invalidCommand{}
	}
	// do not update cards after revealing it
	if p.RoomState == roomShow {
		return nil
	}
	if player, ok := p.Players[cmd.playerID]; ok {
		player.Card = string(cmd.card)
		player.PlayerState = playerSelectedCard
		return nil
	}
	return playerNotFoundError{cmd.playerID}
}

type roomShowCommand struct{}

func (cmd roomShowCommand) Do(ctx context.Context, p *room) error {
	p.RoomState = roomShow
	return nil
}

type roomResetCommand struct{}

func (cmd roomResetCommand) Do(ctx context.Context, p *room) error {
	p.RoomState = roomSelecting
	for _, playerState := range p.Players {
		playerState.Card = ""
		playerState.PlayerState = playerNoSelectedCard
	}
	return nil
}

type noCommand struct{}

func (n noCommand) Do(ctx context.Context, p *room) error {
	fmt.Println("no command associated")
	return nil
}

func commandFactory(ctx context.Context, cmd Command, client *virtualClient) funcCommand {
	if cmd.Action == actionAddPlayer {
		return addPlayerCommand{
			playerName: cmd.Payload["PlayerName"],
			playerID:   cmd.PlayerID,
			roomID:     cmd.RoomID,
			client:     client,
		}
	}
	if cmd.Action == actionPlayerSelectCard {
		return playerSelectCardCommand{playerID: cmd.PlayerID, card: Card(cmd.Payload["Card"])}
	}
	if cmd.Action == actionRoomShowCard {
		return roomShowCommand{}
	}
	if cmd.Action == actionRoomResetCard {
		return roomResetCommand{}
	}
	return noCommand{}
}
