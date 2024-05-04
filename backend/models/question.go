package models

import "encoding/json"

func UnmarshalQuestion(data []byte) (Vote, error) {
	var r Vote
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Vote) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// question
type Vote struct {
	// 发布频道                   
	Channel         string    `json:"channel"`
	// 评论区                    
	Comments        []Comment `json:"comments"`
	// 推荐费用                   
	Fee             *int64    `json:"fee"`
	// 问题ID                   
	ID              string    `json:"id"`
	// 邀请好友ID                 
	InviteIDS       []string  `json:"invite_ids,omitempty"`
	// 问题选项                   
	Options         []string  `json:"options"`
	// 发起人                    
	Sponsor         string    `json:"sponsor"`
	// 发起时间                   
	Time            string    `json:"time"`
	// 问题题目                   
	Title           string    `json:"title"`
	// 问题描述
	Description     string    `json:"description"`
	// 问题类型：公开/私密/匿名          
	Type            Type      `json:"type"`
	// 投票者和他的建议               
	Voters          []string  `json:"voters"`
}

type Comment struct {
	// 建议         
	Advise string `json:"advise"`
	// 立场         
	Choose string `json:"choose"`
	// 用户名        
	Userid string `json:"userid"`
}

// 问题类型：公开/私密/匿名
type Type string

const (
	Anonymous Type = "anonymous"
	Public    Type = "public"
	Secret    Type = "secret"
)