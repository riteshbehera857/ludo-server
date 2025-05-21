package pawn

import (
	"encoding/json"
	"messaging/common"
)

// PawnMoveMessage is a struct that represents a message that is sent when a pawn is moved.
type PawnMoveMessage struct {
	common.Message
	eventName string
	quadrant  string
	pawn      string
	steps     int
}

// NewPawnMoveMessage creates a new PawnMoveMessage.
func NewPawnMoveMessage(eventName string, quadrant string, pawn string, steps int) *PawnMoveMessage {
	return &PawnMoveMessage{
		eventName: eventName,
		quadrant:  quadrant,
		pawn:      pawn,
		steps:     steps,
	}
}

// GetPawnMoveMessage returns the PawnMoveMessage.
func (m *PawnMoveMessage) GetPawnMoveMessage() PawnMoveMessage {
	return PawnMoveMessage{
		pawn:      m.pawn,
		steps:     m.steps,
		quadrant:  m.quadrant,
		eventName: m.eventName,
	}
}

// GetEventName returns the event name of the PawnMoveMessage.
func (m *PawnMoveMessage) GetEventName() string {
	return m.eventName
}

// GetQuadrant returns the quadrant of the PawnMoveMessage.
func (m *PawnMoveMessage) GetQuadrant() string {
	return m.quadrant
}

// GetPawn returns the pawn of the PawnMoveMessage.
func (m *PawnMoveMessage) GetPawn() string {
	return m.pawn
}

// GetSteps returns the steps of the PawnMoveMessage.
func (m *PawnMoveMessage) GetSteps() int {
	return m.steps
}

// ToJSON converts the PawnMoveMessage to a JSON string.
func (m *PawnMoveMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
		Quadrant  string `json:"quadrant"`
		Pawn      string `json:"pawn"`
		Steps     int    `json:"steps"`
	}{
		EventName: m.eventName,
		Quadrant:  m.quadrant,
		Pawn:      m.pawn,
		Steps:     m.steps,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToObject converts a JSON string to a PawnMoveMessage.
func (m *PawnMoveMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
		Quadrant  string `json:"quadrant"`
		Pawn      string `json:"pawn"`
		Steps     int    `json:"steps"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &PawnMoveMessage{}, err
	}

	return &PawnMoveMessage{
		eventName: intermediate.EventName,
		quadrant:  intermediate.Quadrant,
		pawn:      intermediate.Pawn,
		steps:     intermediate.Steps,
	}, nil
}
