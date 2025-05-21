package board

import (
	"encoding/json"
	"messaging/common"
)

// GameEndMessage is a message that is sent to the clients when the game is over and a winner has been determined.
type GameEndMessage struct {
	common.Message
	eventName     string
	winner        string
	winningAmount int
	responseCode  int
}

// NewGameWinnerMessage creates a new GameEndMessage.
func NewGameEndMessage(eventName string, winner string, winningAmount int, responseCode int) *GameEndMessage {
	return &GameEndMessage{
		eventName:     eventName,
		winner:        winner,
		winningAmount: winningAmount,
		responseCode:  responseCode,
	}
}

// GetGameEndMessage returns a copy of the GameEndMessage.
func (m *GameEndMessage) GetGameEndMessage() GameEndMessage {
	return GameEndMessage{
		winner:       m.winner,
		responseCode: m.responseCode,
	}
}

// ToJSON returns the JSON representation of the GameEndMessage.
func (m *GameEndMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName     string `json:"eventName"`
		Winner        string `json:"winner"`
		WinningAmount int    `json:"winningAmount"`
		ResponseCode  int    `json:"responseCode"`
	}{
		EventName:     m.eventName,
		Winner:        m.winner,
		WinningAmount: m.winningAmount,
		ResponseCode:  m.responseCode,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *GameEndMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName     string `json:"eventName"`
		Winner        string `json:"winner"`
		WinningAmount int    `json:"winningAmount"`
		ResponseCode  int    `json:"responseCode"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &GameEndMessage{}, err
	}

	return &GameEndMessage{
		eventName:     intermediate.EventName,
		winner:        intermediate.Winner,
		winningAmount: intermediate.WinningAmount,
		responseCode:  intermediate.ResponseCode,
	}, nil
}
