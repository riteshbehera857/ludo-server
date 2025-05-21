package board

import (
	"encoding/json"
	"messaging/common"
)

// DisconnectionMessage is a struct that represents a message that is sent when a player disconnects.
type DisconnectionMessage struct {
	common.Message
	eventName string
	player    string
}

// NewDisconnectionMessage creates a new DisconnectionMessage.
func NewDisconnectionMessage(eventName string, player string) *DisconnectionMessage {
	return &DisconnectionMessage{
		eventName: eventName,
		player:    player,
	}
}

// GetDisconnectionMessage returns the DisconnectionMessage.
func (m *DisconnectionMessage) GetDisconnectionMessage() DisconnectionMessage {
	return DisconnectionMessage{
		player:    m.player,
		eventName: m.eventName,
	}
}

// GetEventName returns the event name of the DisconnectionMessage.
func (m *DisconnectionMessage) GetEventName() string {

	return m.eventName
}

// GetPlayer returns the player of the DisconnectionMessage.
func (m *DisconnectionMessage) GetPlayer() string {

	return m.player
}

// ToJSON converts the DisconnectionMessage to a JSON string.
func (m *DisconnectionMessage) ToJSON() (string, error) {
	jsonData, err := json.Marshal(&struct {
		EventName string `json:"eventName"`
		Player    string `json:"player"`
	}{
		EventName: m.eventName,
		Player:    m.player,
	})
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

// ToObject converts a JSON string to a DisconnectionMessage.
func (m *DisconnectionMessage) ToObject(data string) (common.Message, error) {
	var intermediate struct {
		EventName string `json:"eventName"`
		Player    string `json:"player"`
	}

	err := json.Unmarshal([]byte(data), &intermediate)

	if err != nil {
		return &DisconnectionMessage{}, err
	}

	return NewDisconnectionMessage(intermediate.EventName, intermediate.Player), nil
}
