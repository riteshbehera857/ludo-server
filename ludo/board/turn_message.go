package board

import (
	"encoding/json"
	"ludo/pawn"
	"messaging/common"
)

type TurnMessage struct {
	eventName string
	turn      string
	positions []pawn.PawnPositions
}

func NewTurnMessage(eventName string, Turn string, pawnPositions []pawn.PawnPositions) *TurnMessage {
	return &TurnMessage{
		turn:      Turn,
		eventName: eventName,
		positions: pawnPositions,
	}
}

func (m *TurnMessage) GetTurnMessage() TurnMessage {
	return TurnMessage{
		turn: m.turn,
	}
}

func (m *TurnMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		Turn      string               `json:"turn"`
		EventName string               `json:"eventName"`
		Positions []pawn.PawnPositions `json:"positions"`
	}{
		Turn:      m.turn,
		EventName: m.eventName,
		Positions: m.positions,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *TurnMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		Turn      string               `json:"turn"`
		EventName string               `json:"eventName"`
		Positions []pawn.PawnPositions `json:"positions"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &TurnMessage{}, err
	}

	return &TurnMessage{
		turn:      intermediate.Turn,
		eventName: intermediate.EventName,
		positions: intermediate.Positions,
	}, nil
}
