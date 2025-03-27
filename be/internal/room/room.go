package room

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type room struct {
	users  map[string]*user
	events chan baseMessage
	done   chan struct{}
}

func newRoom() *room {
	return &room{
		users:  make(map[string]*user),
		events: make(chan baseMessage),
		done:   make(chan struct{}),
	}
}

func (r *room) run() {
	for {
		select {
		case event := <-r.events:
			r.handleEvent(event)
		case <-r.done:
			return
		}
	}
}

func (r *room) stop() {
	close(r.done)
}

func (r *room) join(u *user) {
	r.events <- baseMessage{
		Type: eventJoin,
		From: u.innerID,
		Data: u,
	}
	r.listen(u.userID, u.innerID, u.conn)
}

func (r *room) listen(userID int64, innerID string, conn *websocket.Conn) {
	defer func() {
		r.events <- baseMessage{Type: eventLeave, From: innerID}
	}()

	messageChan := make(chan []byte)
	errChan := make(chan error)

	go func() {
		defer close(messageChan)
		defer close(errChan)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			messageChan <- message
		}
	}()

	ticker := time.NewTicker(websocketPingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("room_listen: ping user %d(%s): %v", userID, innerID, err)
				return
			}
		case input := <-messageChan:
			var message baseMessage
			if err := json.Unmarshal(input, &message); err != nil {
				log.Printf("room_listen: parse message from user %d(%s): %v", userID, innerID, err)
				continue
			}
			message.From = innerID // override the from field (prevent spoofing)
			r.events <- message
		case err := <-errChan:
			log.Printf("room_listen: read from user %d(%s): %v", userID, innerID, err)
			return
		}
	}
}

func (r *room) handleEvent(event baseMessage) {
	switch event.Type {
	case eventJoin:
		u := event.Data.(*user)
		r.users[u.innerID] = u

		r.sendUserJoined(event.From, u.name, u.avatar)
	case eventLeave:
		_ = r.users[event.From].conn.Close()
		delete(r.users, event.From)

		r.sendUserLeft(event.From)
	case eventChatMessage:
		if msg, ok := unmarshalClientData[chatMessage](event.Data); ok {
			r.sendChatMessage(event.From, msg)
		}
	case eventOffer, eventAnswer, eventIceCandidate:
		if msg, ok := unmarshalClientData[webRTCMessage](event.Data); ok {
			r.forwardMessage(event, msg.To)
		}
	}
}

func (r *room) sendUserJoined(innerID string, name, avatar string) {
	for _, u := range r.users {
		// send info message to the user who joined only
		if u.innerID == innerID {
			msg := baseMessage{
				Type: eventInfo,
				From: innerID,
				Data: infoMessage{Users: createInfoUsers(r.users, innerID)},
			}
			u.send(&msg)
			continue
		}

		msg := baseMessage{
			Type: eventJoin,
			From: innerID,
			Data: joinMessage{Name: name, Avatar: avatar},
		}
		u.send(&msg)
	}
}

func (r *room) sendUserLeft(innerID string) {
	for _, u := range r.users {
		msg := baseMessage{
			Type: eventLeave,
			From: innerID,
		}
		u.send(&msg)
	}
}

func (r *room) sendChatMessage(innerID string, chatMsg chatMessage) {
	for _, u := range r.users {
		msg := baseMessage{
			Type: eventChatMessage,
			From: innerID,
			Data: chatMsg,
		}
		u.send(&msg)
	}
}

func (r *room) forwardMessage(msg baseMessage, to string) {
	if targetUser, ok := r.users[to]; ok {
		targetUser.send(&msg)
	}
}

func createInfoUsers(users map[string]*user, exclude string) []infoUser {
	infoUsers := make([]infoUser, 0, len(users)-1)
	for _, u := range users {
		if u.innerID == exclude {
			continue
		}
		infoUsers = append(infoUsers, infoUser{InnerID: u.innerID, Name: u.name, Avatar: u.avatar})
	}
	return infoUsers
}

func unmarshalClientData[T any](input any) (T, bool) {
	data, ok := input.(string)
	if !ok {
		log.Printf("unmarshal: unexpected data type: %T", input)
		return *new(T), false
	}

	var dst T
	if err := json.Unmarshal([]byte(data), &dst); err != nil {
		log.Printf("unmarshal: decode data: %v", err)
		return *new(T), false
	}

	return dst, true
}
