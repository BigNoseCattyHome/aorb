package models

type User struct {
	// 与数据库交互的实体，需修改
	ID       uint32 `gorm:"not null; primary_key; autoIncrement"`               // 用户id
	UserName string `gorm:"not null; unique; size: 32; index" redis:"UserName"` // 用户名
	Password string `gorm:"not null" redis:"Password"`                          // 密码
	Role     int    `gorm:"default:1" redis:"Role"`                             // 角色(普通用户或者管理员)
	Coins    int32  `gorm:"default:0" redis:"Coins"`                            // 金币数量
	Avatar   string `redis:"Avatar"`                                            // 头像地址
}

type UserRequest struct {
	UserId  uint32 `json:"user_id" binding:"required"`
	ActorId uint32 `json:"actor_id" binding:"required"`
}

type UserResponse struct {
	StatusCode int32  `json:"status_code"` // 状态码，0-成功，其他值-失败
	StatusMsg  string `json:"status_msg"`  // 返回状态描述
}
