/*
 * aorb
 *
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * API version: 1.0.0
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package models

type User struct {

	// 用户ID
	Id string `json:"id"`

	// 用户昵称
	Nickname string `json:"nickname"`

	// 用户头像
	Avatar string `json:"avatar"`

	// 用户的金币数
	Coins float32 `json:"coins"`

	// 用户金币流水记录
	CoinsRecord *[]string `json:"coins_record"`

	// 关注者
	Followed []string `json:"followed"`

	// 被关注者
	Follower []string `json:"follower"`

	// 屏蔽好友
	Blacklist []string `json:"blacklist"`

	// 发起过的问题
	QuestionsAsk []string `json:"questions_ask"`

	// 回答过的问题
	QuestionsAsw []string `json:"questions_asw"`

	// 关注的频道
	Channels []string `json:"channels"`
}
