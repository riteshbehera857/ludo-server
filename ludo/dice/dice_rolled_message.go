package dice

import (
	"encoding/json"
	"messaging/common"
)

type DiceRolledMessage struct {
	common.Message
	eventName    string
	number       int
	quadrant     string
	movablePawns []string
}

func NewDiceRolledMessage(
	eventName string,
	number int,
	quadrant string,
	movablePawns []string,
) *DiceRolledMessage {
	return &DiceRolledMessage{
		number:       number,
		eventName:    eventName,
		quadrant:     quadrant,
		movablePawns: movablePawns,
	}
}

func (m *DiceRolledMessage) GetDiceRolledMessage() DiceRolledMessage {

	return DiceRolledMessage{
		number:       m.number,
		quadrant:     m.quadrant,
		movablePawns: m.movablePawns,
	}
}

func (m *DiceRolledMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		Number       int      `json:"number"`
		EventName    string   `json:"eventName"`
		MovablePawns []string `json:"movablePawns"`
		Quadrant     string   `json:"quadrant"`
	}{
		Number:       m.number,
		EventName:    m.eventName,
		MovablePawns: m.movablePawns,
		Quadrant:     m.quadrant,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *DiceRolledMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		Number       int      `json:"number"`
		EventName    string   `json:"eventName"`
		Quadrant     string   `json:"quadrant"`
		MovablePawns []string `json:"movablePawns"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return nil, err
	}

	return &DiceRolledMessage{
		number:       intermediate.Number,
		eventName:    intermediate.EventName,
		quadrant:     intermediate.Quadrant,
		movablePawns: intermediate.MovablePawns,
	}, nil
}
