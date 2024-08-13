package models

import "time"

type Message struct {
	// 与数据库交互的实体
	MessageUuid  string    `json:"message_uuid" bson:"messageUuid,omitempty"`
	FromUserName string    `json:"from_user_name" bson:"fromUserName,omitempty"`
	ToUserName   string    `json:"to_user_name" bson:"toUserName,omitempty"`
	Content      string    `json:"content" bson:"content,omitempty"`
	HasBeenRead  bool      `json:"has_been_read" bson:"hasBeenRead"`
	MessageType  string    `json:"message_type" bson:"messageType,omitempty"`
	CreateAt     time.Time `json:"create_at" bson:"createAt,omitempty"`
}
