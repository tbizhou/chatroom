package websocket

import (
	"log"
	"sync"
)

// 全局Hub实例
var globalHub *Hub
var hubOnce sync.Once

// GetHub 获取全局Hub单例
func GetHub() *Hub {
	hubOnce.Do(func() {
		globalHub = NewHub()
	})
	return globalHub
}

// Hub 连接管理中心，维护所有活跃的客户端连接
type Hub struct {
	Rooms map[string]*Room

	// 消息存储
	messageStore MessageInterface

	// 并发安全锁
	mu sync.RWMutex

	// 停止信号
	stopChan chan struct{}
	done     chan struct{}
}

// SetMessageStore 设置消息存储
func (h *Hub) SetMessageStore(store MessageInterface) {
	h.messageStore = store
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		Rooms:    make(map[string]*Room),
		stopChan: make(chan struct{}),
		done:     make(chan struct{}),
	}
}

func (h *Hub) GetRoom(roomId string) *Room {
	h.mu.Lock()
	defer h.mu.Unlock()
	if room, ok := h.Rooms[roomId]; ok {
		return room
	}
	room := NewRoom(roomId, h.messageStore)
	h.Rooms[roomId] = room
	go room.run()
	log.Printf("Created new Room: %s", roomId)
	return room
}
