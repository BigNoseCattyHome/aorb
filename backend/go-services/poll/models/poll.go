package models

import (
	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	voteModels "github.com/BigNoseCattyHome/aorb/backend/go-services/vote/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"time"
)

type Poll struct {
	PollUuid     string                  `json:"poll_uuid" bson:"pollUuid, omitempty"`
	Title        string                  `json:"title" bson:"title, omitempty"`
	Option1      string                  `json:"option1" bson:"option1, omitempty"`
	Option2      string                  `json:"option2" bson:"option2, omitempty"`
	OptionsCount []uint32                `json:"options_count" bson:"optionsCount, omitempty"`
	PollType     string                  `json:"poll_type" bson:"pollType, omitempty"`
	UserName     string                  `json:"username" bson:"userName, omitempty"`
	CommentList  []commentModels.Comment `json:"comment_list" bson:"commentList, omitempty"`
	VoteList     []voteModels.Vote       `json:"vote_list" bson:"voteList, omitempty"`
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
