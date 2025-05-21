package quadrant

import (
	"encoding/json"
	"messaging/common"
)

type QuadrantSelectMessage struct {
	common.Message
	eventName string
	quadrant  string
}

func NewQuadrantSelectMessage(eventName string, quadrant string) *QuadrantSelectMessage {
	return &QuadrantSelectMessage{
		eventName: eventName,
		quadrant:  quadrant,
	}
}

func (m *QuadrantSelectMessage) GetQuadrant() string {
	return m.quadrant
}

func (m *QuadrantSelectMessage) GetQuadrantSelectMessage() QuadrantSelectMessage {
	return QuadrantSelectMessage{
		quadrant: m.quadrant,
	}
}

func (m *QuadrantSelectMessage) GetEventName() string {
	return m.eventName
}

func (m *QuadrantSelectMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
		Quadrant  string `json:"quadrant"`
	}{
		EventName: m.eventName,
		Quadrant:  m.quadrant,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *QuadrantSelectMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
		Quadrant  string `json:"quadrant"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &QuadrantSelectMessage{}, err
	}

	return &QuadrantSelectMessage{
		eventName: m.eventName,
		quadrant:  intermediate.Quadrant,
	}, nil
}
