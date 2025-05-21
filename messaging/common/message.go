package common

import (
	"encoding/json"
	"fmt"
)

type Message interface {
	ToJSON() (string, error)
	ToObject(string) (Message, error)
}

type SocketMessage struct {
	Message
	eventName    string
	errorCode    int
	errorMessage string
}

func (sm SocketMessage) GetEventName() string {
	return sm.eventName
}

func NewSocketMessage(eventName string, errorCode int, errorMessage string) SocketMessage {
	return SocketMessage{
		eventName:    eventName,
		errorCode:    errorCode,
		errorMessage: errorMessage,
	}
}

func (sm SocketMessage) ToJSON() (string, error) {
	resp, err := json.Marshal(map[string]interface{}{
		"eventName":    sm.eventName,
		"errorCode":    sm.errorCode,
		"errorMessage": sm.errorMessage,
	})
	if err != nil {
		return "", err
	}
	return string(resp), nil
}

func (sm SocketMessage) ToObject(msg string) (Message, error) {
	var msgMap map[string]interface{}
	err := json.Unmarshal([]byte(msg), &msgMap)
	if err != nil {
		return SocketMessage{}, fmt.Errorf("failed to unmarshal message: %v", err)
	}

	// Check if eventName exists and is string
	eventName, ok := msgMap["eventName"]
	if !ok {
		return SocketMessage{}, fmt.Errorf("missing required field: eventName")
	}

	eventNameStr, ok := eventName.(string)
	if !ok {
		return SocketMessage{}, fmt.Errorf("eventName must be a string")
	}

	// Get optional fields with defaults
	errorCode := 0
	if code, exists := msgMap["errorCode"]; exists {
		if codeFloat, ok := code.(float64); ok {
			errorCode = int(codeFloat)
		}
	}

	errorMessage := ""
	if msg, exists := msgMap["errorMessage"]; exists {
		if msgStr, ok := msg.(string); ok {
			errorMessage = msgStr
		}
	}

	return NewSocketMessage(eventNameStr, errorCode, errorMessage), nil
}
