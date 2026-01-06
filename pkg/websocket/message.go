package websocket

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        string `json:"id"`
	RoomID    string `json:"room_id"`
	SenderID  string `json:"sender_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

func NewMessage(roomId, senderId, content string) *Message {
	return &Message{
		ID:        uuid.New().String(),
		RoomID:    roomId,
		SenderID:  senderId,
		Content:   content,
		Timestamp: time.Now().UnixMilli(),
	}
}

func (m *Message) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func ParseMessage(data []byte) (*Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return &msg, err
}
