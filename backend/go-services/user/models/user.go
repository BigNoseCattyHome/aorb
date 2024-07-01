package models

import "time"

type User struct {
	ID       uint32    `json:"id" bson:"_id,omitempty"`
	Username string    `json:"username" bson:"username,omitempty"`
	Nickname string    `json:"nickname" bson:"nickname,omitempty"`
	Password string    `json:"password" bson:"password,omitempty"`
	Avatar   string    `json:"avatar" bson:"avatar,omitempty"`
	CreateAt time.Time `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at" bson:"update_at,omitempty"`
	DeleteAt time.Time `json:"delete_at" bson:"delete_at,omitempty"`
}

type UserReq struct {
	UserId  uint32 `form:"user_id" binding:"required"`
	ActorId uint32 `form:"actor_id" binding:"required"`
}

type UserRes struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}
