package room

import "time"

const (
	websocketPingInterval = 30 * time.Second
)

type eventType string

const (
	// server events

	eventJoin  eventType = "join"
	eventLeave eventType = "leave"
	eventInfo  eventType = "info"

	// client events

	eventChatMessage  eventType = "chat-message"
	eventOffer        eventType = "offer"
	eventAnswer       eventType = "answer"
	eventIceCandidate eventType = "ice-candidate"
)

type (
	baseMessage struct {
		Type eventType `json:"type"`
		From string    `json:"from"`
		Data any       `json:"data,omitempty"`
	}

	joinMessage struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	infoUser struct {
		InnerID string `json:"inner_id"`
		Name    string `json:"name"`
		Avatar  string `json:"avatar"`
	}

	infoMessage struct {
		Users []infoUser `json:"users"`
	}

	chatMessage struct {
		Content string    `json:"content"`
		Ts      time.Time `json:"ts"`
	}

	webRTCMessage struct {
		To      string `json:"to"`
		Content string `json:"content"`
	}
)
