package board

import (
	"encoding/json"
	"messaging/common"
)

type GameStartMessage struct {
	common.Message
	eventName string
}

func NewGameStartMessage(eventName string) *GameStartMessage {
	return &GameStartMessage{
		eventName: eventName,
	}
}

func (m *GameStartMessage) GetGameStartMessage() GameStartMessage {
	return GameStartMessage{}
}

func (m *GameStartMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
	}{
		EventName: m.eventName,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *GameStartMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &GameStartMessage{}, err
	}

	return &GameStartMessage{
		eventName: intermediate.EventName,
	}, nil
}
