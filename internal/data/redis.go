package data

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chatroom/pkg/websocket"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(client *redis.Client) *RedisClient {
	return &RedisClient{client: client}
}

func (s *RedisClient) SaveMessage(ctx context.Context, msg *websocket.Message) error {
	if s.client == nil {
		return nil // Redis未配置，跳过保存
	}
	msgKey := fmt.Sprintf("chat:message:%s", msg.ID)
	msgData, _ := json.Marshal(msg)

	if err := s.client.Set(ctx, msgKey, msgData, 0).Err(); err != nil {
		return err
	}

	// 2. 将消息ID添加到房间消息列表
	listKey := fmt.Sprintf("chat:room:%s:messages", msg.RoomID)
	return s.client.LPush(ctx, listKey, msg.ID).Err()
}

// GetRoomMessages 获取房间历史消息
func (s *RedisClient) GetRoomMessages(ctx context.Context, roomID string, limit int64) ([]*websocket.Message, error) {
	if s.client == nil {
		return []*websocket.Message{}, nil // Redis未配置，返回空
	}
	listKey := fmt.Sprintf("chat:room:%s:messages", roomID)

	msgIDs, err := s.client.LRange(ctx, listKey, 0, limit-1).Result()
	if err != nil {
		return nil, err
	}

	messages := make([]*websocket.Message, 0, len(msgIDs))
	for _, id := range msgIDs {
		msgKey := fmt.Sprintf("chat:message:%s", id)
		data, err := s.client.Get(ctx, msgKey).Bytes()
		if err != nil {
			continue
		}

		var msg websocket.Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}
		messages = append(messages, &msg)
	}

	return messages, nil
}

func (c *RedisClient) GetRooms(ctx context.Context) []string {
	// 从房间集合中获取所有房间ID
	rooms, err := c.client.SMembers(ctx, "chat:rooms").Result()
	if err != nil {
		return []string{}
	}
	return rooms
}

func (c *RedisClient) CreatRoom(ctx context.Context, roomId string) error {
	// 将房间ID添加到房间集合中
	return c.client.SAdd(ctx, "chat:rooms", roomId).Err()
}
