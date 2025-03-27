package room

import (
	"encoding/json"
	"log"

	"github.com/escalopa/peer-cast/internal/domain"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type user struct {
	innerID string
	userID  int64
	name    string
	avatar  string
	conn    *websocket.Conn
}

func newUser(u *domain.User, conn *websocket.Conn) *user {
	return &user{
		innerID: uuid.NewString(),
		userID:  u.UserID,
		name:    u.Name,
		avatar:  u.Avatar,
		conn:    conn,
	}
}

func (u *user) send(msg *baseMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("user: marshal message: %v", err)
	}

	err = u.conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		log.Printf("user: send message: %v", err)
	}
}
