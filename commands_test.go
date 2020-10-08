package main

import (
	"testing"
)

func TestAddPlayerCommand(t *testing.T) {
	socketCommand := command{
		action:        "AddPlayer",
		currentPlayer: "worming",
	}
	addPlayer := commandFactory(socketCommand)
	h := hub{}
	addPlayer.Do(&h)

	if len(h.party.players) != 1 {
		t.Errorf("expected to register a player. Expected 1 player, got %d", len(h.party.players))
	}
}
