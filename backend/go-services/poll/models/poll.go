package models

import (
	commentModels "github.com/BigNoseCattyHome/aorb/backend/go-services/comment/models"
	"github.com/BigNoseCattyHome/aorb/backend/go-services/vote/models"
	"github.com/BigNoseCattyHome/aorb/backend/rpc/poll"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Poll struct {
	PollUuid     string                  `json:"poll_uuid" bson:"pollUuid, omitempty"`
	UserName     string                  `json:"username" bson:"userName, omitempty"`
	PollType     string                  `json:"poll_type" bson:"pollType, omitempty"`
	Title        string                  `json:"title" bson:"title, omitempty"`
	Options      []string                `json:"options" bson:"options, omitempty"`
	OptionsCount []int64                 `json:"options_count" bson:"optionsCount, omitempty"`
	CommentList  []commentModels.Comment `json:"comment_list" bson:"commentList, omitempty"`
	VoteList     []models.Vote           `json:"vote_list" bson:"voteList, omitempty"`
	CreateAt     *timestamppb.Timestamp  `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt     *timestamppb.Timestamp  `json:"update_at" bson:"update_at,omitempty"`
	DeleteAt     *timestamppb.Timestamp  `json:"delete_at" bson:"delete_at,omitempty"`
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
