package board

import (
	"encoding/json"
	"messaging/common"
)

// GameWinnerMessage is a message that is sent to the clients when the game is over and a winner has been determined.
type GameWinnerMessage struct {
	common.Message
	eventName     string
	winner        string
	winningAmount int
	responseCode  int
}

// NewGameWinnerMessage creates a new GameWinnerMessage.
func NewGameWinnerMessage(eventName string, winner string, winningAmount int, responseCode int) *GameWinnerMessage {
	return &GameWinnerMessage{
		eventName:     eventName,
		winner:        winner,
		winningAmount: winningAmount,
		responseCode:  responseCode,
	}
}

// GetGameWinnerMessage returns a copy of the GameWinnerMessage.
func (m *GameWinnerMessage) GetGameWinnerMessage() GameWinnerMessage {
	return GameWinnerMessage{
		winner:       m.winner,
		responseCode: m.responseCode,
	}
}

// ToJSON returns the JSON representation of the GameWinnerMessage.
func (m *GameWinnerMessage) ToJSON() (string, error) {
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

func (m *GameWinnerMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName     string `json:"eventName"`
		Winner        string `json:"winner"`
		WinningAmount int    `json:"winningAmount"`
		ResponseCode  int    `json:"responseCode"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &GameWinnerMessage{}, err
	}

	return &GameWinnerMessage{
		eventName:     intermediate.EventName,
		winner:        intermediate.Winner,
		winningAmount: intermediate.WinningAmount,
		responseCode:  intermediate.ResponseCode,
	}, nil
}
