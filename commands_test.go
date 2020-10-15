package main

import (
	"context"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestAddPlayerCommand(t *testing.T) {
	ctx := newTestContext()
	playerID := 123
	roomID := 456
	playerName := "worming"
	socketCommand := Command{
		Action:   actionAddPlayer,
		PlayerID: playerID,
	}

	addPlayer := commandFactory(ctx, socketCommand, &virtualClient{})
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	addPlayer.Do(ctx, p)

	if len(p.Players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(p.Players))
	}
}

func TestPlayerSelectCardCommand(t *testing.T) {
	ctx := newTestContext()
	expectedCardValue := "3"
	playerID := 123
	roomID := 456
	playerName := "worming"

	payload := make(map[string]string)
	payload["Card"] = expectedCardValue
	socketCommand := Command{
		Action:   actionPlayerSelectCard,
		PlayerID: playerID,
		Payload:  payload,
	}
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	applyCommand(ctx, p, socketCommand)

	if len(p.Players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(p.Players))
		return
	}
	wormingPlayer := p.Players[playerID]
	if wormingPlayer.PlayerState != playerSelectedCard {
		t.Errorf("expected player to have selected card state. Got %v", wormingPlayer.PlayerState)
	}

	if wormingPlayer.Card != expectedCardValue {
		t.Errorf("expected player to have selected card %s. Got %v", expectedCardValue, wormingPlayer.PlayerState)
	}

}

func TestPlayerSelectCardInvalidCard(t *testing.T) {
	ctx := newTestContext()
	expectedCardValue := "inexisting"
	playerName := "worming"
	playerID := 123
	roomID := 456

	payload := make(map[string]string)
	payload["Card"] = expectedCardValue
	socketCommand := Command{
		Action:   actionPlayerSelectCard,
		PlayerID: playerID,
		Payload:  payload,
	}
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	applyCommand(ctx, p, socketCommand)

	if len(p.Players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(p.Players))
		return
	}

	wormingPlayer := p.Players[playerID]

	if wormingPlayer.PlayerState != playerNoSelectedCard {
		t.Errorf("expected player to have selected card state. Got %v", wormingPlayer.PlayerState)
	}

	if wormingPlayer.Card == expectedCardValue {
		t.Errorf("expected player to not have selected card %s as value is not permitted. Got %v", expectedCardValue, wormingPlayer.Card)
	}
}

func TestResetCard(t *testing.T) {
	ctx := newTestContext()
	playerName := "worming"
	playerID := 123
	roomID := 456
	socketCommand := Command{
		Action:   actionRoomResetCard,
		PlayerID: playerID,
	}
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	p.Players[playerID].Card = "should be reset"
	p.Players[playerID].PlayerState = playerSelectedCard
	applyCommand(ctx, p, socketCommand)

	if len(p.Players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(p.Players))
	}

	wormingPlayer := p.Players[playerID]

	if wormingPlayer.Card != "" {
		t.Errorf("expected player to not have selected empty card. Got %v", wormingPlayer.Card)
	}
	if wormingPlayer.PlayerState != playerNoSelectedCard {
		t.Errorf("expected player to have no selected card state. Got %v", wormingPlayer.PlayerState)
	}
}

func TestSelectCardDuringRevealIsDoingNothing(t *testing.T) {
	ctx := newTestContext()
	playerName := "worming"
	playerID := 123
	roomID := 456
	socketCommand := Command{
		Action:   actionRoomShowCard,
		PlayerID: playerID,
	}
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	p.Players[playerID].Card = "3"
	applyCommand(ctx, p, socketCommand)

	payload := make(map[string]string)
	payload["Card"] = "5"
	socketCommandSelectCard := Command{
		Action:   actionPlayerSelectCard,
		PlayerID: playerID,
		Payload:  payload,
	}
	applyCommand(ctx, p, socketCommandSelectCard)

	wormingPlayer := p.Players[playerID]

	if wormingPlayer.Card != "3" {
		t.Errorf("expected previous card 3, got %v", wormingPlayer.Card)
	}
}

func newSinglePlayerHub(playerID, roomID int, playerName string) (*hub, *room) {
	h := newHub(conf{})
	room := newRoom(roomID)
	h.room[roomID] = room
	room.Players[playerID] = &PlayerSelectionState{playerName, playerNoSelectedCard, ""}
	return h, room
}

func applyCommand(ctx context.Context, p *room, socketCommand Command) {
	cmd := commandFactory(ctx, socketCommand, &virtualClient{})
	cmd.Do(ctx, p)
}

func newTestContext() context.Context {
	logger := logrus.NewEntry(logrus.New())
	return context.WithValue(context.Background(), loggerKey, logger)
}
