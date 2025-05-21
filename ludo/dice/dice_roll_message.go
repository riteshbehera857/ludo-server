package dice

import (
	"encoding/json"
	"messaging/common"
)

type DiceRollMessage struct {
	common.Message
	eventName string
}

func NewDiceRollMessage(eventName string) *DiceRollMessage {
	return &DiceRollMessage{
		eventName: eventName,
	}
}

func (m *DiceRollMessage) GetDiceRollMessage() DiceRollMessage {
	return DiceRollMessage{
		eventName: m.eventName,
	}
}

func (m *DiceRollMessage) ToJSON() (string, error) {
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

func (m *DiceRollMessage) ToObject(data string) (DiceRollMessage, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return DiceRollMessage{}, err
	}

	return DiceRollMessage{
		eventName: intermediate.EventName,
	}, nil
}
