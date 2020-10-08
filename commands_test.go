package main

import (
	"testing"
)

func TestAddPlayerCommand(t *testing.T) {
	playerID := 123
	roomID := 456
	playerName := "worming"
	socketCommand := Command{
		Action:   actionAddPlayer,
		PlayerID: playerID,
	}

	addPlayer := commandFactory(socketCommand, &virtualClient{})
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	addPlayer.Do(p)

	if len(p.Players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(p.Players))
	}
}

func TestPlayerSelectCardCommand(t *testing.T) {
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
	applyCommand(p, socketCommand)

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
	applyCommand(p, socketCommand)

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
	applyCommand(p, socketCommand)

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
	playerName := "worming"
	playerID := 123
	roomID := 456
	socketCommand := Command{
		Action:   actionRoomShowCard,
		PlayerID: playerID,
	}
	_, p := newSinglePlayerHub(playerID, roomID, playerName)
	p.Players[playerID].Card = "3"
	applyCommand(p, socketCommand)

	payload := make(map[string]string)
	payload["Card"] = "5"
	socketCommandSelectCard := Command{
		Action:   actionPlayerSelectCard,
		PlayerID: playerID,
		Payload:  payload,
	}
	applyCommand(p, socketCommandSelectCard)

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

func applyCommand(p *room, socketCommand Command) {
	cmd := commandFactory(socketCommand, &virtualClient{})
	cmd.Do(p)
}
