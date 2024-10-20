package models

import "time"

type Vote struct {
	// 与mongoDB交互的实体
	VoteUuid     string    `json:"vote_uuid" bson:"voteUuid,omitempty"`
	VoteUserName string    `json:"vote_username" bson:"voteUserName,omitempty"`
	Choice       string    `json:"choice" bson:"choice,omitempty"`
	CreateAt     time.Time `json:"create_at" bson:"createAt,omitempty"`
}
