package board

import (
	"encoding/json"
	"log"
	"messaging/common"
)

// BoardWaitingPlayersMessage is a message that is sent to the clients when the game is waiting for players to join.
type BoardWaitingPlayersMessage struct {
	common.Message
	eventName               string
	waitingPlayers          []Player
	newPlayer               Player
	playerSelectingQuadrant Player
}

// NewBoardWaitingPlayersMessage creates a new BoardWaitingPlayersMessage.
func NewBoardWaitingPlayersMessage(eventName string, waitingPlayers []Player, newPlayer Player, playerSelectingQuadrant Player) *BoardWaitingPlayersMessage {
	log.Printf("Creating new BoardWaitingPlayersMessage with eventName: %s, waitingPlayers: %v, newPlayer: %v, playerSelectingQuadrant: %v", eventName, waitingPlayers, newPlayer, playerSelectingQuadrant)
	return &BoardWaitingPlayersMessage{
		eventName:               eventName,
		waitingPlayers:          waitingPlayers,
		newPlayer:               newPlayer,
		playerSelectingQuadrant: playerSelectingQuadrant,
	}
}

// GetBoardWaitingPlayersMessage returns a copy of the BoardWaitingPlayersMessage.
func (m *BoardWaitingPlayersMessage) GetBoardWaitingPlayersMessage() BoardWaitingPlayersMessage {
	return BoardWaitingPlayersMessage{}
}

// ToJSON returns the JSON representation of the BoardWaitingPlayersMessage.
func (m *BoardWaitingPlayersMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName               string   `json:"eventName"`
		WaitingPlayers          []Player `json:"waitingPlayers"`
		NewPlayer               Player   `json:"newPlayer"`
		PlayerSelectingQuadrant Player   `json:"playerSelectingQuadrant"`
	}{
		EventName:               m.eventName,
		WaitingPlayers:          m.waitingPlayers,
		NewPlayer:               m.newPlayer,
		PlayerSelectingQuadrant: m.playerSelectingQuadrant,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *BoardWaitingPlayersMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName               string   `json:"eventName"`
		WaitingPlayers          []Player `json:"waitingPlayers"`
		NewPlayer               Player   `json:"newPlayer"`
		PlayerSelectingQuadrant Player   `json:"playerSelectingQuadrant"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)
	if err != nil {
		return nil, err
	}

	return &BoardWaitingPlayersMessage{
		eventName:               intermediate.EventName,
		waitingPlayers:          intermediate.WaitingPlayers,
		newPlayer:               intermediate.NewPlayer,
		playerSelectingQuadrant: intermediate.PlayerSelectingQuadrant,
	}, nil

}
