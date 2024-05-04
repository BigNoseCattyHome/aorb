//    user, err := UnmarshalUser(bytes)
//    bytes, err = user.Marshal()

package models

import "encoding/json"

func UnmarshalUser(data []byte) (User, error) {
	var r User
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *User) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type User struct {
	// 用户头像
	Avatar string `json:"avatar"`
	// 屏蔽好友
	Blacklist []string `json:"blacklist"`
	// 关注的频道
	Channels []string `json:"channels"`
	// 用户的金币数
	Coins float64 `json:"coins"`
	// 用户金币流水记录
	CoinsRecord []string `json:"coins_record"`
	// 关注者
	Followed []string `json:"followed"`
	// 被关注者
	Follower []string `json:"follower"`
	// 用户ID
	ID string `json:"id"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 发起过的问题
	QuestionsAsk []string `json:"questions_ask"`
	// 回答过的问题
	QuestionsAsw []string `json:"questions_asw"`
}
