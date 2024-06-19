package models

// 一条金币流水记录
type CoinRecord struct {
	// 消耗的金币数
	Consume int64 `json:"consume"`
	// 为其投币的问题ID
	QuestionID string `json:"question_id"`
	// 使用者的ID
	UserID string `json:"user_id"`
}
