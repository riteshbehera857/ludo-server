package board

import (
	"encoding/json"
	"ludo/ludo_board_constants"
	"ludo/pawn"
	"messaging/common"
)

type BoardReconnectionMessage struct {
	common.Message
	eventName    string
	Participants []ParticipantInfo
	positions    []pawn.PawnPositions
}

func NewBoardReconnectionMessage(participants []ParticipantInfo, positions []pawn.PawnPositions) *BoardReconnectionMessage {
	return &BoardReconnectionMessage{
		eventName:    ludo_board_constants.BOARD_RECONNECTION,
		Participants: participants,
		positions:    positions,
	}
}

func (m *BoardReconnectionMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName    string               `json:"eventName"`
		Participants []ParticipantInfo    `json:"participants"`
		Positions    []pawn.PawnPositions `json:"positions"`
	}{
		EventName:    m.eventName,
		Participants: m.Participants,
		Positions:    m.positions,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
