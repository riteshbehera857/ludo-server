package pawn

import (
	"encoding/json"
	"messaging/common"
)

type ValidationError struct {
	Message         string `json:"message"`
	CurrentLocation string `json:"currentLocation"`
}

// PawnMovedMessage is a struct that represents a message that is sent when a pawn is moved.
type PawnMovedMessage struct {
	common.Message
	eventName        string
	pawn             string
	steps            int
	initialPosition  int
	finalPosition    int
	initialIndex     int
	finalIndex       int
	isAtHome         bool
	quadrant         string
	capturedPawns    []string
	responseCode     int
	validationErrors []ValidationError
	positions        []PawnPositions
}

// NewPawnMovedMessage creates a new PawnMovedMessage.
func NewPawnMovedMessage(data map[string]interface{}) *PawnMovedMessage {
	msg := &PawnMovedMessage{}

	// Only set fields that exist in the map
	if v, ok := data["eventName"].(string); ok {
		msg.eventName = v
	}
	if v, ok := data["pawn"].(string); ok {
		msg.pawn = v
	}
	if v, ok := data["steps"].(int); ok {
		msg.steps = v
	}
	if v, ok := data["initialPosition"].(int); ok {
		msg.initialPosition = v
	}
	if v, ok := data["finalPosition"].(int); ok {
		msg.finalPosition = v
	}
	if v, ok := data["initialIndex"].(int); ok {
		msg.initialIndex = v
	}
	if v, ok := data["finalIndex"].(int); ok {
		msg.finalIndex = v
	}
	if v, ok := data["isAtHome"].(bool); ok {
		msg.isAtHome = v
	}
	if v, ok := data["quadrant"].(string); ok {
		msg.quadrant = v
	}
	if v, ok := data["capturedPawns"].([]string); ok {
		msg.capturedPawns = v
	}
	if v, ok := data["responseCode"].(int); ok {
		msg.responseCode = v
	}
	if v, ok := data["validationErrors"].([]ValidationError); ok {
		if len(v) == 0 || v[0].Message == "" {
			msg.validationErrors = []ValidationError{}
		} else {
			msg.validationErrors = v
		}
	}
	if v, ok := data["positions"].([]PawnPositions); ok {
		msg.positions = v
	}
	return msg
}

// ToJSON converts the PawnMovedMessage to a JSON string.
func (m *PawnMovedMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName        string            `json:"eventName"`
		Pawn             string            `json:"pawn"`
		Steps            int               `json:"steps"`
		Quadrant         string            `json:"quadrant"`
		ResponseCode     int               `json:"responseCode"`
		InitialPosition  int               `json:"initialPosition"`
		FinalPosition    int               `json:"finalPosition"`
		InitialIndex     int               `json:"initialIndex"`
		FinalIndex       int               `json:"finalIndex"`
		IsAtHome         bool              `json:"isAtHome"`
		CapturedPawns    []string          `json:"capturedPawns"`
		ValidationErrors []ValidationError `json:"validationErrors"`
		Postions         []PawnPositions   `json:"positions"`
	}{
		EventName:        m.eventName,
		Pawn:             m.pawn,
		Steps:            m.steps,
		Quadrant:         m.quadrant,
		ResponseCode:     m.responseCode,
		InitialPosition:  m.initialPosition,
		FinalPosition:    m.finalPosition,
		InitialIndex:     m.initialIndex,
		FinalIndex:       m.finalIndex,
		IsAtHome:         m.isAtHome,
		CapturedPawns:    m.capturedPawns,
		ValidationErrors: m.validationErrors,
		Postions:         m.positions,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToObject converts a JSON string to a PawnMovedMessage.
func (m *PawnMovedMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName        string            `json:"eventName"`
		Pawn             string            `json:"pawn"`
		Steps            int               `json:"steps"`
		Quadrant         string            `json:"quadrant"`
		ResponseCode     int               `json:"responseCode"`
		InitialPosition  int               `json:"initialPosition"`
		FinalPosition    int               `json:"finalPosition"`
		InitialIndex     int               `json:"initialIndex"`
		FinalIndex       int               `json:"finalIndex"`
		IsAtHome         bool              `json:"isAtHome"`
		CapturedPawns    []string          `json:"capturedPawns"`
		ValidationErrors []ValidationError `json:"validationErrors"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &PawnMovedMessage{}, err
	}

	return &PawnMovedMessage{
		eventName:        intermediate.EventName,
		pawn:             intermediate.Pawn,
		steps:            intermediate.Steps,
		quadrant:         intermediate.Quadrant,
		responseCode:     intermediate.ResponseCode,
		initialPosition:  intermediate.InitialPosition,
		finalPosition:    intermediate.FinalPosition,
		initialIndex:     intermediate.InitialIndex,
		finalIndex:       intermediate.FinalIndex,
		isAtHome:         intermediate.IsAtHome,
		capturedPawns:    intermediate.CapturedPawns,
		validationErrors: intermediate.ValidationErrors,
	}, nil
}

type PawnPosition struct {
	Name            string `json:"name"`
	CurrentPosition int    `json:"currentPosition"`
}

type PawnPositions struct {
	Quadrant      string         `json:"quadrant"`
	PawnPositions []PawnPosition `json:"pawnPositions"`
}

func NewPawnPosition(name string, currentPosition int) *PawnPosition {
	return &PawnPosition{
		Name:            name,
		CurrentPosition: currentPosition,
	}
}

func (pp *PawnPositions) AddPosition(name string, currentPosition int) {
	pp.PawnPositions = append(pp.PawnPositions, *NewPawnPosition(name, currentPosition))
}
