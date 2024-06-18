// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    responseLogin, err := UnmarshalResponseLogin(bytes)
//    bytes, err = responseLogin.Marshal()

package models

import "encoding/json"

func UnmarshalResponseLogin(data []byte) (ResponseLogin, error) {
	var r ResponseLogin
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ResponseLogin) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ResponseLogin struct {
	// 过期时间
	ExpiresIn int64 `json:"expires_in"`
	// 消息
	Message string `json:"message"`
	// 操作是否成功
	Success bool `json:"success"`
	// JWT令牌
	Token string `json:"token"`
	// 令牌类型
	TokenType string `json:"token_type"`
	// simpleuser
	User SimpleUser `json:"user"`
}
