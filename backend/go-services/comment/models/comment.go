package models

import (
	"github.com/BigNoseCattyHome/aorb/backend/rpc/comment"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Comment struct {
	ID       primitive.ObjectID `json:"_id,omitempty"`
	UserId   primitive.ObjectID `json:"user_id,omitempty"`
	PollId   primitive.ObjectID `json:"poll_id"`
	Content  string             `json:"comment,omitempty"`
	CreateAt time.Time          `json:"create_at,omitempty"`
	DeleteAt time.Time          `json:"delete_at"`
}

type ActionCommentReq struct {
	Token       string `json:"token" binding:"required"`
	ActorId     int    `json:"actor_id"`
	PollId      int    `json:"poll_id"`
	ActionType  int    `json:"action_type"`
	CommentText string `json:"comment_text"`
	CommentId   int    `json:"comment_id"`
}

type ActionCommentRes struct {
	StatusCode int             `json:"status_code"`
	StatusMsg  string          `json:"status_msg"`
	Comment    comment.Comment `json:"comment"`
}

type ListCommentReq struct {
	Token   string `json:"token"`
	ActorId int    `json:"actor_id"`
	PollId  int    `json:"poll_id" binding:"-"`
}

type ListCommentRes struct {
	StatusCode  int                `json:"status_code"`
	StatusMsg   string             `json:"status_msg"`
	CommentList []*comment.Comment `json:"comment_list"`
}

type CountCommentReq struct {
	Token   string `json:"token"`
	ActorId int    `json:"actor_id"`
	PollId  int    `json:"poll_id"`
}

type CountCommentRes struct {
	StatusCode   int    `json:"status_code"`
	StatusMsg    string `json:"status_msg"`
	CommentCount int    `json:"comment_count"`
}
