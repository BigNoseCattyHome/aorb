package models

import (
	"time"

	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
)

type Poll struct {
	PollUuid     string                  `json:"poll_uuid" bson:"pollUuid, omitempty"`
	Title        string                  `json:"title" bson:"title, omitempty"`
	Options      []string                `json:"options" bson:"options, omitempty"`
	OptionsCount []uint32                `json:"options_count" bson:"optionsCount, omitempty"`
	Content      string                  `json:"content" bson:"content, omitempty"`
	PollType     string                  `json:"poll_type" bson:"pollType, omitempty"`
	UserName     string                  `json:"username" bson:"userName, omitempty"`
	CommentList  []Comment `json:"comment_list" bson:"commentList, omitempty"`
	VoteList     []Vote       `json:"vote_list" bson:"voteList, omitempty"`
	CreateAt     time.Time               `json:"create_at" bson:"createAt, omitempty"`
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
