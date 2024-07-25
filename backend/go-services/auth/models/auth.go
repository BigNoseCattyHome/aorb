package models

import (
	"encoding/json"
	"time"
)

// 一条金币流水记录
type CoinRecord struct {
	// 消耗的金币数
	Consume int64 `json:"consume"`
	// 为其投币的问题ID
	QuestionID string `json:"question_id"`
	// 使用者的ID
	UserID string `json:"user_id"`
}

func UnmarshalRequestLogin(data []byte) (LoginRequest, error) {
	var r LoginRequest
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginRequest struct {
	// 设备ID
	DeviceID string `json:"deviceId"`
	// 用户名/用户ID
	Username string `json:"username"`
	// 随机数
	Nonce string `json:"nonce"`
	// 密码的md5摘要
	Password string `json:"password"`
	// 时间戳
	Timestamp time.Time `json:"timestamp"`
}

func UnmarshalResponseLogin(data []byte) (LoginResponse, error) {
	var r LoginResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *LoginResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type LoginResponse struct {
	// 访问令牌的过期时间
	ExpiresIn int64 `json:"expires_in"`
	// 消息
	Message string `json:"message"`
	// 刷新令牌
	RefreshToken string `json:"refresh_token"`
	// 操作是否成功
	Success bool `json:"success"`
	// JWT令牌
	Token     string     `json:"token"`
	TokenType string     `json:"token_type"`
	User      SimpleUser `json:"user"`
}

func UnmarshalResponseVerify(data []byte) (VerifyResponse, error) {
	var r VerifyResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *VerifyResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type VerifyResponse struct {
	// 过期时间
	Exp int64 `json:"exp"`
	// 用户ID
	UserID string `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 是否有效
	Valid bool `json:"valid"`
}
type SimpleUser struct {
	// 头像
	Avatar string `json:"avatar"`
	// 用户ID
	ID string `json:"id"`
	// IP归属地
	Ipaddress string `json:"ipaddress"`
	// 昵称
	Nickname string `json:"nickname"`
}

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
	Username string    `json:"username" bson:"username"`
	CreateAt time.Time `json:"create_at" bson:"create_at,omitempty"`
	UpdateAt time.Time `json:"update_at" bson:"update_at,omitempty"`
	DeleteAt time.Time `json:"delete_at" bson:"delete_at"`
}
