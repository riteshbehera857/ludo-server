package board

import (
	"encoding/json"
	"messaging/common"
)

type Player struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type ParticipantInfo struct {
	Player   Player `json:"player"`
	Quadrant string `json:"quadrant"`
}

type BoardJoinedMessage struct {
	common.Message
	eventName                  string
	participants               []ParticipantInfo
	playerSelectingTheQuadrant Player
}

func NewBoardJoinedMessage(eventName string, participants []ParticipantInfo, playerSelectingTheQuadrant Player) *BoardJoinedMessage {
	return &BoardJoinedMessage{
		eventName:                  eventName,
		participants:               participants,
		playerSelectingTheQuadrant: playerSelectingTheQuadrant,
	}
}

func (m *BoardJoinedMessage) GetBoardJoinedMessage() BoardJoinedMessage {
	return BoardJoinedMessage{
		participants: m.participants,
	}
}

func (m *BoardJoinedMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName                  string            `json:"eventName"`
		Participants               []ParticipantInfo `json:"participants"`
		PlayerSelectingTheQuadrant Player            `json:"playerSelectingTheQuadrant"`
	}{
		EventName:                  m.eventName,
		Participants:               m.participants,
		PlayerSelectingTheQuadrant: m.playerSelectingTheQuadrant,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *BoardJoinedMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName                  string            `json:"eventName"`
		Participants               []ParticipantInfo `json:"participants"`
		PlayerSelectingTheQuadrant Player            `json:"playerSelectingTheQuadrant"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &BoardJoinedMessage{}, err
	}

	return &BoardJoinedMessage{
		eventName:                  intermediate.EventName,
		participants:               intermediate.Participants,
		playerSelectingTheQuadrant: intermediate.PlayerSelectingTheQuadrant,
	}, nil
}
