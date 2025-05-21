package board

import (
	"encoding/json"
	"messaging/common"
)

type BoardSelectingQuadrantMessage struct {
	common.Message
	eventName string
	player    Player
}

func NewBoardSelectingQuadrantMessage(eventName string, player Player) *BoardSelectingQuadrantMessage {
	return &BoardSelectingQuadrantMessage{
		eventName: eventName,
		player:    player,
	}
}

func (m *BoardSelectingQuadrantMessage) GetBoardSelectingQuadrantMessage() BoardSelectingQuadrantMessage {
	return BoardSelectingQuadrantMessage{
		player: m.player,
	}
}

func (m *BoardSelectingQuadrantMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
		Player    Player `json:"player"`
	}{
		EventName: m.eventName,
		Player:    m.player,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *BoardSelectingQuadrantMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
		Player    Player `json:"player"`
	}
	err := json.Unmarshal([]byte(data), &intermediate)
	if err != nil {
		return nil, err
	}

	return &BoardSelectingQuadrantMessage{
		eventName: intermediate.EventName,
		player:    intermediate.Player,
	}, nil
}
