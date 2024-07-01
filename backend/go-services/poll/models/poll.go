package models

import (
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"time"
)

type Poll struct {
	ID          uint32    `json:"id" bson:"_id, omitempty"`
	UserId      uint32    `json:"userId" bson:"userId, omitempty"`
	PollType    string    `json:"poll_type" bson:"poll_type, omitempty"`
	Title       string    `json:"title" bson:"title, omitempty"`
	Options     []string  `json:"options" bson:"options, omitempty"`
	OptionsRate []float64 `json:"optionsRate" bson:"optionsRate, omitempty"`
	CreateAt    time.Time `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt    time.Time `json:"update_at" bson:"update_at,omitempty"`
	DeleteAt    time.Time `json:"delete_at" bson:"delete_at,omitempty"`
}

type ListPollReq struct {
	LatestTime string `form:"latest_time"`
	ActorId    int    `form:"actor_id"`
}

type ListPollRes struct {
	StatusCode int          `json:"status_code"`
	StatusMsg  string       `json:"status_msg"`
	NextTime   *int64       `json:"next_time,omitempty"`
	PollList   []*poll.Poll `json:"poll_list,omitempty"`
}
