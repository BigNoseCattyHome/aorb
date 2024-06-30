package models

import (
	"encoding/json"
)

func UnmarshalUser(data []byte) (User, error) {
	var r User
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// User represents a user in the system.
type User struct {
	// 用户头像
	Avatar string `json:"avatar" bson:"avatar"`
	// 屏蔽好友
	Blacklist []string `json:"blacklist" bson:"blacklist"`
	// 用户的金币数
	Coins float64 `json:"coins" bson:"coins"`
	// 用户金币流水记录
	CoinsRecord []CoinRecord `json:"coins_record" bson:"coins_record"`
	// 关注者
	Followed []string `json:"followed" bson:"followed"`
	// 被关注者
	Follower []string `json:"follower" bson:"follower"`
	// 用户ID，这个是Objectid，由服务端mongodb生成，不支持修改
	ID string `json:"id" bson:"_id"`
	// IP归属地
	Ipaddress string `json:"ipaddress" bson:"ipaddress"`
	// 用户昵称
	Nickname string `json:"nickname" bson:"nickname"`
	// 用户密码
	Password string `json:"password" bson:"password"`
	// 发起过的问题
	QuestionsAsk []string `json:"questions_ask" bson:"questions_ask"`
	// 回答过的问题
	QuestionsAsw []string `json:"questions_asw" bson:"questions_asw"`
	// 收藏的问题
	QuestionsCollect []string `json:"questions_collect" bson:"questions_collect"`
	// 用户登录名
	Username string `json:"username" bson:"username"`
}
