package websocket

import (
	"context"
	"log"
)

type MessageInterface interface {
	SaveMessage(ctx context.Context, msg *Message) error
	GetRoomMessages(ctx context.Context, roomID string, limit int64) ([]*Message, error)
}

type Room struct {
	ID           string
	Clients      map[*Client]bool
	Broadcast    chan *Message
	Register     chan *Client
	Unregister   chan *Client
	messageStore MessageInterface
}

func NewRoom(id string, store MessageInterface) *Room {
	return &Room{
		ID:           id,
		Clients:      make(map[*Client]bool),
		Broadcast:    make(chan *Message),
		Register:     make(chan *Client),
		Unregister:   make(chan *Client),
		messageStore: store,
	}
}

func (r *Room) run() {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
			log.Printf("Client joined Room: %s. Total: %d", r.ID, len(r.Clients))
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				close(client.send)
				log.Printf("Client left Room: %s. Total: %d", r.ID, len(r.Clients))
			}
		case message := <-r.Broadcast:
			// 保存消息到存储
			if r.messageStore != nil {
				if err := r.messageStore.SaveMessage(context.Background(), message); err != nil {
					log.Printf("Failed to save message: %v", err)
				}
			}
			// 广播给所有客户端
			for client := range r.Clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(r.Clients, client)
				}
			}
		}
	}
}
