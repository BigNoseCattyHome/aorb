package models

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
