package board

import (
	"encoding/json"
	"messaging/common"
)

type BoardBetFailedMessage struct {
	common.Message
	eventName string
	message   string
}

func NewBoardBetFailedMessage(eventName string, message string) *BoardBetFailedMessage {
	return &BoardBetFailedMessage{
		eventName: eventName,
		message:   message,
	}
}

func (m *BoardBetFailedMessage) GetBoardBetFailedMessage() BoardBetFailedMessage {
	return BoardBetFailedMessage{
		eventName: m.eventName,
		message:   m.message,
	}
}

func (m *BoardBetFailedMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
		Message   string `json:"message"`
	}{
		EventName: m.eventName,
		Message:   m.message,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *BoardBetFailedMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
		Message   string `json:"message"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &BoardBetFailedMessage{}, err
	}

	return &BoardBetFailedMessage{
		eventName: intermediate.EventName,
		message:   intermediate.Message,
	}, nil
}
