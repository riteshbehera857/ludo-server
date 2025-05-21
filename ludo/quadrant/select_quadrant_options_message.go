package quadrant

import (
	"encoding/json"
	"messaging/common"
)

type SelectQuadrantMessage struct {
	common.SocketMessage
	eventName    string
	quadrants    []string
	responseCode int
}

func NewSelectQuadrantMessage(eventName string, quadrants []string, responseCode int) *SelectQuadrantMessage {
	return &SelectQuadrantMessage{
		eventName:    eventName,
		quadrants:    quadrants,
		responseCode: responseCode,
	}
}

func (m *SelectQuadrantMessage) GetSelectQuadrantMessage() SelectQuadrantMessage {
	return SelectQuadrantMessage{
		quadrants:    m.quadrants,
		responseCode: m.responseCode,
	}
}

func (m *SelectQuadrantMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName    string   `json:"eventName"`
		Quadrants    []string `json:"quadrants"`
		ResponseCode int      `json:"responseCode"`
	}{
		EventName:    m.eventName,
		Quadrants:    m.quadrants,
		ResponseCode: m.responseCode,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *SelectQuadrantMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName    string   `json:"eventName"`
		Quadrants    []string `json:"quadrants"`
		ResponseCode int      `json:"responseCode"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &SelectQuadrantMessage{}, err
	}

	return &SelectQuadrantMessage{
		eventName:    intermediate.EventName,
		quadrants:    intermediate.Quadrants,
		responseCode: intermediate.ResponseCode,
	}, nil
}
