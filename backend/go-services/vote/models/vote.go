package models

import "google.golang.org/protobuf/types/known/timestamppb"

type Vote struct {
	// 与mongoDB交互的实体
	VoteUuid     string                 `json:"vote_uuid" bson:"voteUuid,omitempty"`
	VoteUserName string                 `json:"vote_username" bson:"voteUserName,omitempty"`
	Choice       string                 `json:"choice" bson:"choice,omitempty"`
	CreateAt     *timestamppb.Timestamp `json:"create_at" bson:"create_at,omitempty"`
	DeleteAt     *timestamppb.Timestamp `json:"delete_at" bson:"delete_at"`
}
