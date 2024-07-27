package models

import (
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Comment struct {
	// 与mongodb交互的Comment实体
	CommentUuid      string                 `json:"commentUuid" bson:"commentUuid,omitempty"`
	ReviewerUserName string                 `json:"reviewerUserName" bson:"reviewerUserName,omitempty"`
	Content          string                 `json:"content" bson:"content,omitempty"`
	CreateAt         *timestamppb.Timestamp `json:"create_at" bson:"create_at,omitempty"`
	DeleteAt         *timestamppb.Timestamp `json:"delete_at" bson:"delete_at"`
}

type ActionCommentReq struct {
	Token       string `form:"token" binding:"required"`
	ActorId     int    `form:"actor_id"`
	PollId      int    `form:"poll_id"`
	ActionType  int    `form:"action_type"`
	CommentText string `form:"comment_text"`
	CommentId   int    `form:"comment_id"`
}

type ActionCommentRes struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	Comment    comment.Comment `json:"comment"`
}

type ListCommentReq struct {
	Token   string `form:"token"`
	ActorId int    `form:"actor_id"`
	PollId  int    `form:"poll_id" binding:"-"`
}

type ListCommentRes struct {
	StatusCode  int                `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	CommentList []*comment.Comment `json:"comment_list"`
}

type CountCommentReq struct {
	Token   string `form:"token"`
	ActorId int    `form:"actor_id"`
	PollId  int    `form:"poll_id"`
}

type CountCommentRes struct {
	StatusCode   int    `json:"status_code"`
	StatusMsg    string `json:"status_msg"`
	CommentCount int    `json:"comment_count"`
}
