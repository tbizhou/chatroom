package websocket

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// writeWait 写入超时时间
	writeWait = 10 * time.Second

	// pongWait 等待 Pong 响应的超时时间
	pongWait = 60 * time.Second

	// pingPeriod 发送 Ping 的间隔 (必须小于 pongWait)
	pingPeriod = (pongWait * 9) / 10

	// maxMessageSize 最大消息大小 (64KB)
	maxMessageSize = 64 * 1024

	// sendBufferSize 发送缓冲区大小
	sendBufferSize = 256
)

// Client WebSocket 客户端连接
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	room   *Room
	userID string
	send   chan *Message
}

// NewClient 创建新的客户端
func NewClient(hub *Hub, conn *websocket.Conn, room *Room, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		room:   room,
		userID: userID,
		send:   make(chan *Message, sendBufferSize),
	}
}

// ReadPump 读取客户端消息
// 在独立的 goroutine 中运行，从 WebSocket 连接读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("[WS Client] User %s: read error: %v", c.userID, err)
			}
			break
		}

		// 处理收到的消息
		msg := NewMessage(c.room.ID, c.userID, string(message))
		c.room.Broadcast <- msg
	}
}

// WritePump 向客户端发送消息
// 在独立的 goroutine 中运行，向 WebSocket 连接发送消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub 关闭了 send 通道
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			data, _ := message.ToJSON()
			w.Write(data)

			// 批量发送队列中的消息
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				msg := <-c.send
				msgData, _ := msg.ToJSON()
				w.Write(msgData)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
