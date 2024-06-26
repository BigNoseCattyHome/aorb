package models

import (
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"time"
)

type Comment struct {
	// 与mongodb交互的Comment实体
	ID       string    `json:"id" bson:"_id,omitempty"`
	UserId   string    `json:"user_id" bson:"user_id,omitempty"`
	PollId   string    `json:"poll_id" bson:"poll_id,omitempty"`
	Content  string    `json:"content" bson:"comment,omitempty"`
	CreateAt time.Time `json:"create_at" bson:"create_at,omitempty"`
	DeleteAt time.Time `json:"delete_at" bson:"delete_at"`
}

type ActionCommentReq struct {
	Token       string `json:"token" binding:"required"`
	ActorId     string `json:"actor_id"`
	PollId      string `json:"poll_id"`
	ActionType  int    `json:"action_type"`
	CommentText string `json:"comment_text"`
	CommentId   string `json:"comment_id"`
}

type ActionCommentRes struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	Comment    comment.Comment `json:"comment"`
}

type ListCommentReq struct {
	Token   string `json:"token"`
	ActorId string `json:"actor_id"`
	PollId  string `json:"poll_id" binding:"-"`
}

type ListCommentRes struct {
	StatusCode  int                `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	CommentList []*comment.Comment `json:"comment_list"`
}

type CountCommentReq struct {
	Token   string `json:"token"`
	ActorId string `json:"actor_id"`
	PollId  string `json:"poll_id"`
}

type CountCommentRes struct {
	StatusCode   int    `json:"status_code"`
	StatusMsg    string `json:"status_msg"`
	CommentCount int    `json:"comment_count"`
}
