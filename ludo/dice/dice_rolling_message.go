package dice

import (
	"encoding/json"
	"messaging/common"
)

type DiceRollingMessage struct {
	common.Message
	eventName string
}

func NewDiceRollingMessage(
	eventName string,
) *DiceRolledMessage {
	return &DiceRolledMessage{
		eventName: eventName,
	}
}

func (m *DiceRollingMessage) ToJSON() (string, error) {
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

func (m *DiceRollingMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return nil, err
	}

	return &DiceRolledMessage{
		eventName: intermediate.EventName,
	}, nil
}
