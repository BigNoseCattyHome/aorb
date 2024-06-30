// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    responseVerify, err := UnmarshalResponseVerify(bytes)
//    bytes, err = responseVerify.Marshal()

package models

import "encoding/json"

func UnmarshalResponseVerify(data []byte) (ResponseVerify, error) {
	var r ResponseVerify
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ResponseVerify) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ResponseVerify struct {
	// 过期时间
	Exp int64 `json:"exp"`
	// 用户ID
	UserID string `json:"user_id"`
	// 用户名
	Username string `json:"username"`
	// 是否有效
	Valid bool `json:"valid"`
}
