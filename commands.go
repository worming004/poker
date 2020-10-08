package main

import (
	"log"
)

type action string

type command struct {
	currentPlayer string
	action        action
}

type funcCommand interface {
	Do(*hub) error
}

type addPlayerCommand struct {
	playerName string
}

func (cmd addPlayerCommand) Do(h *hub) error {
	h.party.players = append(h.party.players, playerSelectionState{player: player{name: cmd.playerName}})
	log.Println("cmd received")
	return nil
}

func commandFactory(cmd command) funcCommand {
	return addPlayerCommand{"worming"}
}
