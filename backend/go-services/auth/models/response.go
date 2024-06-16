// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    response, err := UnmarshalResponse(bytes)
//    bytes, err = response.Marshal()

package models

import "encoding/json"

// 这个生成的文件就包含了返回结构体的两个struct，还有他们的与json数据结构的转换方法

// 这个是从json转换为Response的方法，用于处理接收的数据
func UnmarshalResponse(data []byte) (Response, error) {
	var r Response
	err := json.Unmarshal(data, &r)
	return r, err
}

// 这个是从Response转换为json的方法，用于返回数据
func (r *Response) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Response struct {
	// 消息
	Message string `json:"message"`
	// 操作是否成功
	Success bool `json:"success"`
	// JWT令牌
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UserResponse struct {
	// 头像
	Avatar string `json:"avatar"`
	// 用户ID
	ID string `json:"id"`
	// IP归属地
	Ipaddress string `json:"ipaddress"`
	// 昵称
	Nickname string `json:"nickname"`
}
