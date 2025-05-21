package board

import (
	"encoding/json"
	"messaging/common"
)

type Quadrant struct {
	Name  string   `json:"name"`
	Color string   `json:"color"`
	Pawns []string `json:"pawns"`
	Path  []int    `json:"path"`
}

type GameInitializeMessage struct {
	common.Message
	eventName                  string
	safePositions              []int
	quadrants                  []Quadrant
	autoPlay                   bool
	autoPlayTimer              int
	ticketAmount               int
	playersRequiredToStartGame int
	playerSelectingTheQuadrant Player
}

func NewGameInitializeMessage(eventName string, safePositions []int, quadrants []Quadrant, autoPlay bool, playersRequiredToStartGame int, ticketAmount int, playerSelectingTheQuadrant Player, autoPlayTimer int) *GameInitializeMessage {
	return &GameInitializeMessage{
		eventName:                  eventName,
		safePositions:              safePositions,
		quadrants:                  quadrants,
		autoPlay:                   autoPlay,
		autoPlayTimer:              autoPlayTimer,
		ticketAmount:               ticketAmount,
		playersRequiredToStartGame: playersRequiredToStartGame,
		playerSelectingTheQuadrant: playerSelectingTheQuadrant,
	}
}

func (m *GameInitializeMessage) GetGameInitializeMessage() GameInitializeMessage {
	return GameInitializeMessage{
		eventName:                  m.eventName,
		safePositions:              m.safePositions,
		quadrants:                  m.quadrants,
		autoPlay:                   m.autoPlay,
		ticketAmount:               m.ticketAmount,
		playersRequiredToStartGame: m.playersRequiredToStartGame,
	}
}

func (m *GameInitializeMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName                  string     `json:"eventName"`
		SafePositions              []int      `json:"safePositions"`
		Quadrants                  []Quadrant `json:"quadrants"`
		AutoPlay                   bool       `json:"autoPlay"`
		AutoPlayTimer              int        `json:"autoPlayTimer"`
		PlayersRequiredToStartGame int        `json:"playersRequiredToStartGame"`
		TicketAmount               int        `json:"ticketAmount"`
		PlayerSelectingTheQuadrant Player     `json:"playerSelectingTheQuadrant"`
	}{
		EventName:                  m.eventName,
		SafePositions:              m.safePositions,
		Quadrants:                  m.quadrants,
		AutoPlay:                   m.autoPlay,
		AutoPlayTimer:              m.autoPlayTimer,
		TicketAmount:               m.ticketAmount,
		PlayersRequiredToStartGame: m.playersRequiredToStartGame,
		PlayerSelectingTheQuadrant: m.playerSelectingTheQuadrant,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (m *GameInitializeMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName                  string     `json:"eventName"`
		SafePositions              []int      `json:"safePositions"`
		Quadrants                  []Quadrant `json:"quadrants"`
		AutoPlay                   bool       `json:"autoPlay"`
		AutoPlayTimer              int        `json:"autoPlayTimer"`
		PlayersRequiredToStartGame int        `json:"playersRequiredToStartGame"`
		TicketAmount               int        `json:"ticketAmount"`
		PlayerSelectingTheQuadrant Player     `json:"playerSelectingTheQuadrant"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &GameInitializeMessage{}, err
	}

	return &GameInitializeMessage{
		eventName:                  intermediate.EventName,
		safePositions:              intermediate.SafePositions,
		quadrants:                  intermediate.Quadrants,
		autoPlay:                   intermediate.AutoPlay,
		autoPlayTimer:              intermediate.AutoPlayTimer,
		playersRequiredToStartGame: intermediate.PlayersRequiredToStartGame,
		ticketAmount:               intermediate.TicketAmount,
		playerSelectingTheQuadrant: intermediate.PlayerSelectingTheQuadrant,
	}, nil
}
