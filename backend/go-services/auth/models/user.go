// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    user, err := UnmarshalUser(bytes)
//    bytes, err = user.Marshal()

package models

import "encoding/json"

func UnmarshalUser(data []byte) (SimpleUser, error) {
	var r SimpleUser
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *SimpleUser) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// user
type User struct {
	// 用户头像
	Avatar string `json:"avatar"`
	// 屏蔽好友
	Blacklist []string `json:"blacklist"`
	// 用户的金币数
	Coins float64 `json:"coins"`
	// 用户金币流水记录
	CoinsRecord []CoinRecord `json:"coins_record"`
	// 关注者
	Followed []string `json:"followed"`
	// 被关注者
	Follower []string `json:"follower"`
	// 用户ID，这个是Objectid，由服务端mongodb生成，不支持修改
	ID string `json:"id"`
	// IP归属地
	Ipaddress string `json:"ipaddress"`
	// 用户昵称
	Nickname string `json:"nickname"`
	// 用户密码
	Password string `json:"password"`
	// 发起过的问题
	QuestionsAsk []string `json:"questions_ask"`
	// 回答过的问题
	QuestionsAsw []string `json:"questions_asw"`
	// 收藏的问题
	QuestionsCollect []string `json:"questions_collect"`
	// 用户登录名
	Username string `json:"username"`
}
